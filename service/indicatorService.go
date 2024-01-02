package service

import (
	"crypto-trading-bot-go/core"
	"os"
	"os/signal"
	"syscall"
)

func IndicatorMode() {
	LoadHistoryKline()
	go WsTickService()

	symbolsString := ""
	for _, symbol := range core.Config.Trading.Symbols {
		if symbolsString == "" {
			symbolsString = symbol
		} else {
			symbolsString += ", " + symbol
		}
	}

	SendSlack("Indicator Started.\nTracking: " + symbolsString)
	Logger.Warn("Indicator Started")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	SendSlack("Indicator Terminated")
	Logger.Warn("Indicator Terminated")
}
