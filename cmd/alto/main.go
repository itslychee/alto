package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ItsLychee/alto/dsl"
	"github.com/dhowden/tag"
	"github.com/pkg/errors"
)

var SupportedFormats = []tag.FileType{
	tag.MP3,
	tag.M4A,
	tag.M4B,
	tag.M4P,
	tag.ALAC,
	tag.FLAC,
	tag.OGG,
	tag.DSF,
}

type Config struct {
	Path        string `json:"path"`
	Destination string `json:"destination"`
	Source      string `json:"source"`
}

func filepathFunc(dst *string) func(s string) error {
	return func(s string) error {
		if v, err := filepath.Abs(s); err != nil {
			return err
		} else {
			*dst = v
		}
		return nil
	}
}

func main() {
	var config Config
	base, _ := os.UserConfigDir()
	buf, _ := os.ReadFile(filepath.Join(base, "alto", "config.json"))
	json.Unmarshal(buf, &config)
	fmt.Println(config)
	if v, err := filepath.Abs(config.Destination); err == nil {
		config.Destination = v
	}
	if v, err := filepath.Abs(config.Source); err == nil {
		config.Source = v
	}

	flag.Func("config", "custom path to configuration file", func(s string) error {
		buf, err := os.ReadFile(s)
		if err != nil {
			return err
		}
		config = Config{}
		return json.Unmarshal(buf, &config)
	})
	flag.Func("path", "formatting syntax alto should use for files", func(s string) error {
		config.Path = s
		return nil
	})
	flag.Func("source", "where alto should read and index from", filepathFunc(&config.Source))
	flag.Func("destination", "where alto should write to", filepathFunc(&config.Destination))
	flag.Parse()

	if config.Destination == "" || config.Path == "" {
		log.Panicln("path and/or destination must not be nil")
	}

	var sourceIndex []string
	err := filepath.WalkDir(config.Source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return err
		}
		for _, ext := range SupportedFormats {
			if strings.HasSuffix(strings.ToLower(path), strings.ToLower(string(ext))) {
				log.Printf("[%d] indexed: %s", len(sourceIndex)+1, path)
				sourceIndex = append(sourceIndex, path)
				return nil
			}
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
	}
	scope, nodes, err := ParseFormatString(config.Path)
	if err != nil {
		log.Panicln(errors.Wrap(err, "could not compile nodes for provided path"))
	}
	if len(nodes) == 0 {
		log.Panicln("no nodes were sent")
	}

	if err := os.MkdirAll(config.Destination, 0); err != nil {
		log.Panic(err)
	}

	if err := os.Chdir(config.Destination); err != nil {
		log.Panic(err)
	}

index_iter:
	for index, path := range sourceIndex {
		prelimInfo := fmt.Sprintf("[%d/%d]", index+1, len(sourceIndex))
		log.Println(prelimInfo, "opening:", path)
		f, err := os.Open(path)
		if err != nil {
			log.Panicln(errors.Wrap(err, fmt.Sprintf("error while opening %s", path)))
		}

		scope.Functions = map[string]dsl.ASTFunction{}
		scope.Variables = map[string]string{}

		metadata, err := tag.ReadFrom(f)
		if err != nil {
			log.Println(prelimInfo, errors.Wrap(err, "metadata may not be present, error"))
		} else {
			discCurrent, discTotal := metadata.Disc()
			trackCurrent, trackTotal := metadata.Track()
			scope.Variables = map[string]string{
				"trackcurrent": strconv.Itoa(trackCurrent),
				"tracktotal":   strconv.Itoa(trackTotal),
				"disccurrent":  strconv.Itoa(discCurrent),
				"disctotal":    strconv.Itoa(discTotal),
				"year":         strconv.Itoa(metadata.Year()),
				"comment":      metadata.Comment(),
				"format":       string(metadata.Format()),
				"composer":     metadata.Composer(),
				"genre":        metadata.Genre(),
				"albumartist":  metadata.AlbumArtist(),
				"album":        metadata.Album(),
				"artist":       metadata.Artist(),
				"title":        metadata.Title(),
				"filetype":     strings.ToLower(string(metadata.FileType())),
			}
		}
		f.Seek(0, 0)

		rawFilename := strings.Split(filepath.Base(path), ".")

		scope.Variables["filename"] = strings.Join(rawFilename[:len(rawFilename)-1], ".")
		scope.Variables["alto_index"] = strconv.Itoa(index)
		scope.Variables["alto_source"] = config.Source
		scope.Variables["alto_dest"] = config.Destination

		scope.Functions = dsl.DefaultFunctions
		for k, v := range AltoFunctions {
			scope.Functions[k] = v
		}

		var output strings.Builder
		for _, v := range nodes {
			s, err := v.Execute(scope)
			if err != nil {
				if err == ErrSkip {
					log.Println(prelimInfo, "<skip> called")
					continue index_iter
				}
				panic(err)
			}
			output.WriteString(s)
		}
		if output.String() == "" {
			panic("no output string")
		}

		filename := filepath.Join(config.Destination, output.String())
		if err := os.MkdirAll(filepath.Dir(filename), os.ModeDir); err != nil {
			panic(err)
		}

		destFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0)
		if err != nil {
			panic(err)
		}
		// TODO:
		// Allow customization of how alto should handle:
		// * empty output strings
		// * already existing files

		written, err := io.Copy(destFile, f)
		if err != nil {
			panic(err)
		}

		// log.Println(prelimInfo, )
		log.Println(prelimInfo, "finished copying to", filename)
		log.Println(prelimInfo, fmt.Sprintf("results: wrote %d MBs (%d bytes)", written/1000000, written))

	}

}
