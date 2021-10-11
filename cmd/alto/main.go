package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	logging "log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ItsLychee/alto/dsl"
	"github.com/dhowden/tag"
)

var log = logging.New(os.Stderr, "] ", logging.Lmsgprefix|logging.LstdFlags)
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

type Filepath struct {
	val *string
}

func (f Filepath) String() string {
	if f.val != nil {
		return *f.val
	}
	return ""
}

func (f Filepath) Set(value string) error {
	if v, err := filepath.Abs(filepath.Clean(value)); err != nil {
		return err
	} else {
		if value == "" {
			return errors.New("this flag requires a filepath")
		}
		*f.val = v
	}
	return nil
}

func main() {
	var destDirectory string
	var sourceDirectory string = "."
	var format, operation string
	flag.Var(&Filepath{&sourceDirectory}, "source", "directory that alto will read from upon operation")
	flag.Var(&Filepath{&destDirectory}, "destination", "directory that alto will write to upon operation")
	flag.StringVar(&operation, "operation", "copy", "file operation to use, rename and copy are the currently available operations")
	flag.StringVar(&format, "format", "", "format to use to dynamically determine the path upon operation")
	flag.Parse()

	if destDirectory == "" {
		log.Fatalln("-destination must be set to a legible filepath")
	}

	if format == "" {
		log.Fatalln("-format must be set")
	}

	if operation != "copy" && operation != "rename" {
		log.Fatalln("invalid -operation value, expected 'copy' or 'rename'")
	}

	lexer := dsl.NewLexer(format)
	toks, err := lexer.Lex()
	if err != nil {
		log.Fatalln("lexer error:", err)
	}

	parser := dsl.NewParser(toks)
	nodes, err := parser.Parse()
	if err != nil {
		log.Fatalln("parsing error:", err)
	}

	var directoryIndex []string
	err = filepath.WalkDir(sourceDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println("warning:", err)
			return fs.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		for index, val := range SupportedFormats {
			if len(SupportedFormats) == index+1 {
				return nil
			}
			if strings.HasSuffix(strings.ToUpper(path), string(val)) {
				break
			}
		}
		directoryIndex = append(directoryIndex, path)
		log.Printf("[%d] indexed file: %s", len(directoryIndex), filepath.ToSlash(path))
		return nil
	})
	if err != nil {
		panic(err)
	}

	if len(directoryIndex) == 0 {
		log.Fatalln("alto found no supported audio files, ensure your -source directory has music that alto can work with")
	}

	log.Println("starting organization process...")
	if err := os.MkdirAll(destDirectory, 0); err != nil {
		log.Fatalln(err)
	}
	if err := os.Chdir(destDirectory); err != nil {
		log.Fatalln(err)
	}

	var reservedKeywords *regexp.Regexp
	if runtime.GOOS == "windows" {
		reservedKeywords = regexp.MustCompile(`[\pC"*/:<>?\\|]+`)
	} else {
		reservedKeywords = regexp.MustCompile(`[/\x{0}]+`)
	}

	for index, fp := range directoryIndex {
		log.Printf("[%d/%d] traversed on file %s\n", index+1, len(directoryIndex), filepath.ToSlash(fp))
		sourceFile, err := os.Open(fp)
		if err != nil {
			log.Panicln("error while opening file:", fp, err)
		}
		defer sourceFile.Close()

		metadata, err := tag.ReadFrom(sourceFile)
		if err != nil {
			if err == io.EOF || err == tag.ErrNoTagsFound {
				log.Println("warning: could not get any metadata from file")
			} else {
				log.Panicln("error while identifying metadata:", err)
			}
		}
		var scope = dsl.Scope{Variables: make(map[string]string)}
		if metadata != nil {
			trackNumber, trackTotal := metadata.Track()
			discNumber, discTotal := metadata.Disc()
			scope.Variables = map[string]string{
				"title":       reservedKeywords.ReplaceAllString(metadata.Title(), "-"),
				"artist":      reservedKeywords.ReplaceAllString(metadata.Artist(), "-"),
				"album":       reservedKeywords.ReplaceAllString(metadata.Album(), "-"),
				"albumartist": reservedKeywords.ReplaceAllString(metadata.AlbumArtist(), "-"),
				"genre":       reservedKeywords.ReplaceAllString(metadata.Genre(), "-"),
				"composer":    reservedKeywords.ReplaceAllString(metadata.Composer(), "-"),
				"year":        strconv.Itoa(metadata.Year()),

				"tracknumber": strconv.Itoa(trackNumber),
				"tracktotal":  strconv.Itoa(trackTotal),

				"discnumber": strconv.Itoa(discNumber),
				"disctotal":  strconv.Itoa(discTotal),

				"comment":  reservedKeywords.ReplaceAllString(metadata.Comment(), "-"),
				"format":   string(metadata.Format()),
				"filetype": string(metadata.FileType()),
			}
		}

		// Execute AST nodes
		var output strings.Builder
		for _, node := range nodes {
			s, err := node.Execute(scope)
			if err != nil {
				log.Fatalln("error while trying to evaluate the provided format", err)
			}
			output.WriteString(s)
		}

		if output.Len() == 0 {
			log.Println("alto returned an empty path construct, skipping")
			continue
		}

		filename, err := filepath.Abs(output.String())
		if err != nil {
			log.Panicln("could not get an absolute representation of path")
		}

		var extension = filepath.Ext(sourceFile.Name())
		if !strings.HasSuffix(strings.ToLower(filename), extension) {
			filename = strings.Join([]string{filename, extension}, "")
		}

		// Check for path collisions

		var tempFilename = filename
		for counter := 0; ; counter++ {
			_, err := os.Stat(tempFilename)
			if errors.Is(err, os.ErrNotExist) {
				filename = tempFilename
				break
			}
			tempFilename = fmt.Sprintf("%s (%d).%s", filename[:len(filename)-(len(extension)+1)], counter, extension)
			counter++
		}

		if err := os.MkdirAll(filepath.Dir(filename), 0); err != nil {
			log.Panicln("error while creating necessary dirs for dest path", err)
		}

		switch operation {
		case "rename":
			sourceFile.Close()
			if err := os.Rename(sourceFile.Name(), filename); err != nil {
				log.Panicf("error while renaming file\n%s\nto %s", filepath.ToSlash(filename), err)
			}
			log.Printf("[%d/%d] file relocated to %s\n", index+1, len(directoryIndex), filename)

		case "copy":
			sourceFile.Seek(0, 0)

			destFile, err := os.Create(filename)
			if err != nil {
				log.Panicln("error while trying to create dest file", err)
			}
			defer destFile.Close()

			// TODO: perhaps rewrite the copying logic into something
			// more optimized to reduce heavy disk I/O.
			if _, err := io.Copy(destFile, sourceFile); err != nil {
				log.Println("error while copying file", err)
			}

			sourceFile.Close()
			destFile.Close()
			log.Printf("[%d/%d] copied file to %s\n", index+1, len(directoryIndex), filename)
			time.Sleep(50 * time.Millisecond)
		}

	}

}
