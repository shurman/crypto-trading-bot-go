// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/service"
)

//TODO
//record and output filling position

//Download history kline and load
//verify performance

func main() {
	service.LoadConfigs()
	service.InitSlog()

	service.InitStrategyService()

	service.FutureClientInit()
	service.LoadHistoryKline()
	go service.InitWsTickService()

	select {}
}
