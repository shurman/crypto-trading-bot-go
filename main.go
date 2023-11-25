// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//detect fill and exit
//record and output filling position

//Download history kline and load
//verify performance

func main() {
	service.InitStrategyService()

	service.FutureClientInit()
	service.LoadHistoryKline()
	go service.InitWsTickService()

	select {}
}
