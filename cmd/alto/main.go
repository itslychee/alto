package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
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
	Path string `json:"path"`
}

func main() {
	var config Config
	var configPath Filepath
	var destination string
	var source string

	flag.Var(configPath, "config", "custom path to configuration file")
	flag.StringVar(&config.Path, "path", "", "how alto should rename your files")
	flag.StringVar(&source, "source", ".", "where should alto index and read from")
	flag.StringVar(&destination, "destination", ".", "where should alto write to")
	flag.Parse()

	// Configuration loading
	if configPath.String() == "" {
		confdir, _ := os.UserConfigDir()
		defaultConfigFile := filepath.Join(confdir, "alto", "config.json")
		buf, _ := os.ReadFile(defaultConfigFile)
		json.Unmarshal(buf, &config)
	} else {
		buf, err := os.ReadFile(configPath.String())
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(buf, &config); err != nil {
			panic(err)
		}
	}

}
