// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//record and output filling position
//position size for same risk

//Download history kline and load
//verify performance

func main() {
	service.InitStrategyService()

	service.FutureClientInit()
	service.LoadHistoryKline()
	go service.InitWsTickService()

	select {}
}
