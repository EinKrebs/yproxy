package ylogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var Zero = NewZeroLogger("")

func newWriter(filepath string) (*os.File, io.Writer, error) {
	if filepath == "" {
		return nil, os.Stdout, nil
	}
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f, f, nil
}

func NewZeroLogger(filepath string) *zerolog.Logger {
	_, writer, err := newWriter(filepath)
	if err != nil {
		fmt.Printf("FAILED TO INITIALIZED LOGGER: %v", err)
	}
	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}

func UpdateZeroLogLevel(logLevel string) error {
	level := parseLevel(logLevel)
	zeroLogger := Zero.With().Logger().Level(level)
	Zero = &zeroLogger
	return nil
}

func ReloadLogger(filepath string) {
	if filepath == "" { //
		return // this means os.Stdout, so no need to open new file
	}
	newLogger := NewZeroLogger(filepath).Level(Zero.GetLevel())
	Zero = &newLogger
}

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "disabled":
		return zerolog.Disabled
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
