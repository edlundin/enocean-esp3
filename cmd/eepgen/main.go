package main

import (
	"flag"
	"log"

	"github.com/edlundin/enocean-esp3/internal/eepgen"
)

func main() {
	xmlPath := flag.String("xml", "eep268.xml", "path to eep XML")
	outDir := flag.String("out", "pkg/eep/profiles", "output directory")
	flag.Parse()
	if err := eepgen.Generate(*xmlPath, *outDir); err != nil {
		log.Fatal(err)
	}
}
