package common

import (
	"log/slog"
	"path/filepath"
	"strings"
)

// GetJSONHandlerOptions returns customized slog.HandlerOptions for JSON logging.
// Purpose is to provide a consistent, structured JSON log format with simplified source information.
// Example usage:
//
//	logLevel := new(slog.LevelVar)
//	h := slog.NewJSONHandler(os.Stderr, common.GetJSONHandlerOptions(logLevel))
//	slog.SetDefault(slog.New(h))
//
// To use a different log level: `logLevel.Set(slog.LevelDebug)`
// Example log output:
//
//	{
//	  "time": "2024-10-07T09:33:35.310+06:00",
//	  "level": "INFO",
//	  "source": {
//	    "function": "NewServer",
//	    "file": "server.go",
//	    "line": 74
//	  },
//	  "msg": "Swagger Specs available at :8080/swagger/index.html"
//	}
func GetJSONHandlerOptions(logLevel *slog.LevelVar) *slog.HandlerOptions {
	extractFuncName := func(fullFuncName string) string {
		if lastDotIndex := strings.LastIndex(fullFuncName, "."); lastDotIndex != -1 {
			return fullFuncName[lastDotIndex+1:]
		}
		return fullFuncName
	}

	customizeSourceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key != slog.SourceKey {
			return a
		}

		source, ok := a.Value.Any().(*slog.Source)
		if !ok {
			return a
		}

		return slog.Attr{
			Key: slog.SourceKey,
			Value: slog.AnyValue(map[string]any{
				"function": extractFuncName(source.Function),
				"file":     filepath.Base(source.File),
				"line":     source.Line,
			}),
		}
	}

	return &slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel,
		ReplaceAttr: customizeSourceAttr,
	}
}

// GetTextHandlerOptions returns customized slog.HandlerOptions for text-based logging.
// Purpose is to provide a concise text log format with simplified source file information.
// Example usage:
//
//	logLevel := new(slog.LevelVar)
//	h := slog.NewTextHandler(os.Stderr, common.GetTextHandlerOptions(logLevel))
//	slog.SetDefault(slog.New(h))
//
// To use a different log level: `logLevel.Set(slog.LevelDebug)`
// Example log output:
//
//	time=2024-10-07T09:36:51.868+06:00 level=INFO source=server.go:73 msg="Swagger Specs available at :8080/swagger/index.html"
func GetTextHandlerOptions(logLevel *slog.LevelVar) *slog.HandlerOptions {
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
