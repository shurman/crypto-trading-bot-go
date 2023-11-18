// logInitializer
package core

import (
	"log/slog"
)

var (
	Logger *slog.Logger
)

func InitSlog() {
	options := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	Logger = slog.New(ColorHandler(options))
	slog.SetDefault(Logger)
}
