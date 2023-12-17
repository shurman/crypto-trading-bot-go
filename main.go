// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//DTB: find edge failure case / entry rule should base on amplifier
//https://www.ptt.cc/bbs/Trading/M.1538318192.A.FBC.html
//fee calculation
//next strategy
//binance create order

//Future work
//multiple interval(?)

func main() {
	service.InitStrategyService()

	if core.Config.Trading.Mode == "indicator" {
		service.IndicatorMode()
	} else {
		service.BacktestingMode()
	}
}
