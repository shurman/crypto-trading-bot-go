package service

import (
	"crypto-trading-bot-go/core"
	"os"
	"os/signal"
	"syscall"
)

func IndicatorMode() {
	LoadHistoryKline()
	go InitWsTickService()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	if core.Config.Slack.Enable {
		SendSlack("Indicator process terminated!!")
	}

	Logger.Warn("Indicator process terminated!!")
}
