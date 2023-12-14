package service

import (
	"os"
	"os/signal"
	"syscall"
)

func IndicatorMode() {
	LoadHistoryKline()
	go WsTickService()

	SendSlack("Indicator Started")
	Logger.Warn("Indicator Started")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	SendSlack("Indicator Terminated")
	Logger.Warn("Indicator Terminated")
}
