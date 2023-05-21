package main

import (
	"log"
	"os"

	"github.com/mgarlaschelli/logparser/config"
	"github.com/mgarlaschelli/logparser/offset"
	"github.com/mgarlaschelli/logparser/stream"
)

const version string = "1.0.0"

func main() {

	// Read CLI flags

	opts := readCliOptions()

	if !verifyOptions(opts) {
		printUsage()
		os.Exit(1)
	}

	// Parse config file - filters

	filters, err := config.ParseConfigFile(opts.configFile)

	if err != nil {
		log.Fatal(err)
	}

	// Read offset

	offsetMap, err := offset.ReadOffsetFile(opts.offsetFile)

	if err != nil {
		log.Fatal(err)
	}

	// check filters

	err = stream.ParseFiles(opts.logFile, filters, offsetMap, opts.allFiles, opts.ignoreOld)

	if err != nil {
		log.Fatal(err)
	}

	// Remove not existing files from offset map
	newOffsetMap := offset.ClearOffsetMap(offsetMap)

	// Write offset

	offset.WriteOffsetFile(opts.offsetFile, newOffsetMap)
}
