package main

import (
	"flag"
	"io"
	"log"

	"github.com/edlundin/enocean-esp3/internal/eepgen"
)

// main runs the command.
func main() {
	if err := run(flag.CommandLine, nil); err != nil {
		log.Fatal(err)
	}
}

// run executes the command.
func run(fs *flag.FlagSet, args []string) error {
	xmlPath := fs.String("xml", "eep268.xml", "path to eep XML")
	outDir := fs.String("out", "pkg/eep/profiles", "output directory")
	if args != nil {
		fs.SetOutput(io.Discard)
		if err := fs.Parse(args); err != nil {
			return err
		}
	} else {
		flag.Parse()
	}
	return eepgen.Generate(*xmlPath, *outDir)
}
