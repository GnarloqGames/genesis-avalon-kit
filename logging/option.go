package logging

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type Option func() error

func WithKind(newKind string) Option {
	return func() (err error) {
		switch newKind {
		case "json", "text":
			kind = newKind
		default:
			err = fmt.Errorf("invalid kind %s", newKind)
		}

		return
	}
}

func WithOutput(out io.Writer) Option {
	return func() (err error) {
		output = out

		return
	}
}

func WithLevel(newLevel string) Option {
	return func() (err error) {
		switch strings.ToLower(newLevel) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			err = fmt.Errorf("invalid level %s", newLevel)
		}

		return
	}
}
