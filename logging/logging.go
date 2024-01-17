package logging

import "log/slog"

func Logger(options ...Option) (*slog.Logger, error) {
	for _, opt := range options {
		if err := opt(); err != nil {
			return nil, err
		}
	}

	var handler slog.Handler

	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	switch kind {
	case "json":
		handler = slog.NewJSONHandler(output, handlerOpts)
	case "text":
		handler = slog.NewTextHandler(output, handlerOpts)
	}

	return slog.New(handler), nil
}
