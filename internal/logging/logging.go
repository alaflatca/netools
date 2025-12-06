package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Path  string
	Level string
}

func New(cfg Config) *slog.Logger {
	var level slog.Level
	switch strings.ToUpper(cfg.Level) {
	case slog.LevelDebug.String():
		level = slog.LevelDebug
	case slog.LevelInfo.String():
		level = slog.LevelInfo
	case slog.LevelWarn.String():
		level = slog.LevelWarn
	case slog.LevelError.String():
		level = slog.LevelError
	default:

	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if src, ok := a.Value.Any().(slog.Source); ok {
					src.File = filepath.Base(src.File)
					a.Value = slog.AnyValue(src)
				}
			}
			return a
		},
	})
	logger := slog.New(h)

	return logger
}
