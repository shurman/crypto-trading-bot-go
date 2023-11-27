// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//record and output filling position

//Download history kline and load
//verify performance

func main() {
	//service.DownloadRawHistoryKline(service.Config.Trading.Symbol, service.Config.Trading.Interval, 1600000000000, 1500)
	service.InitStrategyService()

	service.LoadHistoryKline()
	go service.InitWsTickService()

	select {}
}
