package main

import (
	"flag"
	"fmt"
	"os"
)

var usage string = fmt.Sprintf("logparser - version %s\nUSAGE: %s -f=<config_file> -o=<offset_file> -l=<log_file> [-a=true] [-i=<ignore_old_files>]", version, os.Args[0])

type Opts struct {
	configFile string
	offsetFile string
	logFile    string
	allFiles   bool
	ignoreOld  int
}

func readCliOptions() Opts {

	opts := Opts{}

	flag.StringVar(&opts.configFile, "f", "", "Configuration file")
	flag.StringVar(&opts.offsetFile, "o", "", "Offset file")
	flag.StringVar(&opts.logFile, "l", "", "Log file")
	flag.BoolVar(&opts.allFiles, "a", false, "Manage all matching files (default: manage only last one)")
	flag.IntVar(&opts.ignoreOld, "i", 3, "Ignore files older than N days")

	flag.Parse()

	return opts
}

func verifyOptions(opts Opts) bool {

	if len(opts.configFile) == 0 {
		return false
	}

	if len(opts.offsetFile) == 0 {
		return false
	}

	if len(opts.logFile) == 0 {
		return false
	}

	return true
}

func printUsage() {
	fmt.Println(usage)
	flag.PrintDefaults()
}
