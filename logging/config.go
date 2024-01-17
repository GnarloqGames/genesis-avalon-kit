package logging

import (
	"io"
	"log/slog"
	"os"
)

var (
	// kind is the output handler kind
	kind string = "json"

	// output is the output the log messages are written to
	output io.Writer = os.Stdout

	// level is the minimum log level to be printed
	level slog.Level = slog.LevelInfo
)
