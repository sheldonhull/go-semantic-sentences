package main

import (

	// "context".
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/peterbourgon/ff/v3"
	"github.com/sheldonhull/go-semantic-sentences/internal/logger"
)

const (
	// exitFail is the exit code if the program
	// fails.
	exitFail   = 1
	MaxSize    = 10
	MaxBackups = 7
	MaxAge     = 7
)

// Logger contains the package level logger provided from internal logger package that wraps up zerolog.
// var log *logger.Logger //nolint: gochecknoglobals

// main configuration from Matt Ryer with minimal logic, passing to run, to allow easier CLI tests.
func main() {
	if err := Run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func Run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New("no arguments")
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	debug := flag.Bool("debug", false, "sets log level to debug and console pretty output")
	source := flag.String("source", "", "source file")
	write := flag.Bool("write", false, "default to stdout, otherwise replace contents of the file")

	// (&debug,
	// 	"debug",
	// 	false,
	// 	"sets log level to debug and console pretty output")

	// ff.WithEnvVarNoPrefix(),

	// ff.WithConfigFileFlag("config"),
	// ff.WithConfigFileParser(fftoml.Parser),
	if err := ff.Parse(flags, args); err != nil {
		return err
	}

	logger.Log.Info().
		Bool("debug", *debug).
		Str("source", *source).
		Bool("write", *write).Msg("flags")

	LogLevel := "info"
	if *debug {
		LogLevel = "debug"
	}

	c := logger.Config{
		Enable:                true,
		ConsoleLoggingEnabled: true,
		EncodeLogsAsJson:      false,
		FileLoggingEnabled:    false,
		Directory:             "",
		Filename:              "",
		MaxSize:               MaxSize,
		MaxBackups:            MaxBackups,
		MaxAge:                MaxAge,
		Level:                 LogLevel,
	}

	_ = logger.InitLogger(c)
	filename := *source

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Log.Error().Err(err).Str("filename", filename).Msg("ReadFile")
		os.Exit(exitFail)
	}

	formatted := FormatSemanticLineBreak(b)
	ioutil.WriteFile(filename, []byte(formatted), os.ModeDevice)

	return nil
}
