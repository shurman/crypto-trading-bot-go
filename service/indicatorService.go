package service

import (
	"os"
	"os/signal"
	"syscall"
)

func IndicatorMode() {
	LoadHistoryKline()
	go WsTickService()

	SendSlack("Indicator process started")
	Logger.Warn("Indicator process started")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	SendSlack("Indicator process terminated")
	Logger.Warn("Indicator process terminated")
}
