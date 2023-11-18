// logHandler
package core

import (
	"log/slog"
	"os"
)

var (
	Logger *slog.Logger
)

func InitSlog() {
	th := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	})

	Logger = slog.New(th)
}
