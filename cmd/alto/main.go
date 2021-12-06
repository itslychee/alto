package main

import (
	"fmt"
	logging "log"
	"os"
	"strings"

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
	nodes, err := ParseFormatString("/one/two/three/four<must(onetwo|three|four)>")
	if err != nil {
		panic(err)
	}

	var scope = dsl.Scope{
		Variables: map[string]string{
			"helloworld": "yo",
		},
		Functions: dsl.DefaultFunctions,
	}

	var builder strings.Builder
	for _, v := range nodes {
		s, err := v.Execute(scope)
		if err != nil {
			panic(err)
		}
		builder.WriteString(s)
	}
	fmt.Println(builder.String())

}
