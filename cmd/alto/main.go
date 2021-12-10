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

func main() {
	var config Config
	base, _ := os.UserConfigDir()
	buf, _ := os.ReadFile(filepath.Join(base, "alto", "config.json"))
	json.Unmarshal(buf, &config)

	flag.Func("config", "custom path to configuration file", func(s string) error {
		buf, err := os.ReadFile(s)
		if err != nil {
			return err
		}
		config = Config{}
		return json.Unmarshal(buf, &config)
	})
	flag.StringVar(&config.Path, "path", "", "formatting syntax alto should use for files")
	flag.StringVar(&config.Source, "source", ".", "where alto should read and index from")
	flag.StringVar(&config.Destination, "destination", "", "where alto should write to")
	flag.Parse()

	if config.Destination == "" || config.Path == "" {
		log.Panicln("path and/or destination must not be nil")
	}

	var sourceIndex []string
	err := filepath.WalkDir(config.Source, func(path string, d fs.DirEntry, err error) error {
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
		return err
	})
	if err != nil {
		log.Panicln(err)
	}
	scope, nodes, err := ParseFormatString(config.Path)
	if err != nil {
		log.Panicln(errors.Wrap(err, "could not compile nodes for provided path"))
	}

	for index, path := range sourceIndex {
		prelimInfo := fmt.Sprintf("[%d/%d]", index+1, len(sourceIndex))
		log.Println(prelimInfo, "opening:", path)
		f, err := os.Open(path)
		if err != nil {
			log.Panicln(errors.Wrap(err, fmt.Sprintf("error while opening %s", path)))
		}
		metadata, err := tag.ReadFrom(f)
		if err != nil {
			log.Panic(errors.Wrap(err, prelimInfo+" could not retrieve metadata"))
		}

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
			"filename":     filepath.Base(path),
			"_index":		strconv.Itoa(index),
		}
		scope.Functions = dsl.DefaultFunctions

		var output strings.Builder
		for _, v := range nodes {
			s, err := v.Execute(scope)
			if err != nil {
				panic(err)
			}
			output.WriteString(s)
		}
		if output.String() == "" {
			panic("no output string")
		}

		f.Seek(0, 0)

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
