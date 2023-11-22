// logInitializerService
package service

import (
	"crypto-trading-bot-go/core"
	"log/slog"
)

var (
	Logger *slog.Logger
)

func InitSlog() {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(Config.System.Loglevel))

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
