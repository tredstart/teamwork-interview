package main

import (
	"flag"
	processor "importer/dataprocessor"
)

type Options struct {
	path    *string
	outFile *string
}

func readOptions() *Options {
	opts := &Options{}
	opts.path = flag.String("path", "./customers.csv", "Path to the file with customer data")
	opts.outFile = flag.String("out", "", "Optional: output file path. If empty program will output results to the terminal")
	flag.Parse()
	return opts
}

func main() {
	opts := readOptions()
	// TODO: improve invalid usage
	processor.ProcessDomainData(*opts.path, *opts.outFile)
}

