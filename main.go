// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//DTB: 2023/11/22 1930 why short / find edge failure case / entry rule should base on amplifier
//next strategy
//binance create order

//Future work
//DTB: state 1,-1 rule rework
//multiple interval(?)

func main() {
	service.InitStrategyService()

	if core.Config.Trading.Mode == "indicator" {
		service.IndicatorMode()
	} else {
		service.BacktestingMode()
	}
}
