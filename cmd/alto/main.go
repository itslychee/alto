package main

import (
	"os"
	logging "log"

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

func main() {
	

}
