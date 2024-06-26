// logInitializerService
package service

import (
	"crypto-trading-bot-go/core"
	"log/slog"
)

var (
	Logger *slog.Logger
)

func init() {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(core.Config.System.Loglevel))

	if err != nil {
		panic(err)
	}

	options := &slog.HandlerOptions{
		Level:     *level,
		AddSource: true,
	}

	Logger = slog.New(core.ColorHandler(options))
	slog.SetDefault(Logger)
}
