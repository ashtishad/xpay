package common

import (
	"log/slog"
	"os"
	"path/filepath"
)

// NewSlogger creates a new text handler with log level,
// Uses custom handler options, strips the full directory path from the source's filename and sets log level.
// Usage:
// logger := common.NewSlogger(slog.LevelDebug)
// slog.SetDefault(logger)
func NewSlogger(logLevel slog.Level) *slog.Logger {
	opts := setHandlerOpts(logLevel)

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

// setHandlerOpts strips the full directory path from the source's filename and sets log level.
func setHandlerOpts(logLevel slog.Level) *slog.HandlerOptions {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			sourceVal, ok := a.Value.Any().(*slog.Source)
			if !ok {
				return a
			}

			sourceVal.File = filepath.Base(sourceVal.File)
		}

		return a
	}

	return &slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel,
		ReplaceAttr: replace,
	}
}
