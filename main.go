// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//DTB phase 2
//try ratio to 1.5 //conclusion not trade btc

//Future work
//state 1,-1 rule rework
//multiple interval(?)

func main() {
	service.InitStrategyService()

	if core.Config.Trading.Mode == "indicator" {
		indicatorMode()
	} else {
		service.BacktestingMode()
	}
}

func indicatorMode() {
	service.LoadHistoryKline()
	go service.InitWsTickService()
	select {}
}
