package main

import (
	"flag"
	processor "importer/dataprocessor"
	"importer/logger"
	"log/slog"
	"os"
)

type Options struct {
	path     *string
	outFile  *string
	logFile  *string
	logLevel *string
}

func readOptions() *Options {
	opts := &Options{}
	opts.path = flag.String("path", "./customers.csv", "Optional: path to the file with customer data")
	opts.outFile = flag.String("out", "", "Optional: output file path. If empty program will output results to the terminal")
	opts.logFile = flag.String("log", "", "Optional: log file path. If empty program will print logs to stderr")
	opts.logLevel = flag.String("log-level", "info", "Optional: sets the level for logging. Options are: debug, info, warn, error. Default is set to info")
	flag.Parse()
	return opts
}

func main() {
	opts := readOptions()

	writer := os.Stderr
	defer func() {
		if writer != os.Stderr {
			writer.Close()
		}
	}()

	if *opts.logFile != "" {
		var err error
		writer, err = os.Create(*opts.logFile)
		if err != nil {
			slog.Warn("Cannot open log file. Redirecting logs to the stderr.", "file_error", err)
		}
	}

	logger.SetupLogger(writer, *opts.logLevel)

	processor.ProcessDomainData(*opts.path, *opts.outFile)
}
