// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//record and output filling position

//verify performance

func main() {
	service.InitStrategyService()

	//onlineMode()
	localMode()
}

func onlineMode() {
	service.LoadHistoryKline()
	go service.InitWsTickService()
	select {}
}

func localMode() {
	//service.DownloadRawHistoryKline(service.Config.Trading.Symbol, service.Config.Trading.Interval, 1500000000000, 1500)
	service.LoadRawHistoryKline(service.Config.Trading.Symbol, service.Config.Trading.Interval)
}
