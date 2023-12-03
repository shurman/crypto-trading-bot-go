// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//multiple coins
//DTB reset after long waiting / state 1,-1 rule rework / phase 2

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
