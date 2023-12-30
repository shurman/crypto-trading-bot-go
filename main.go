// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//入場後被震盪掃出場: 找近幾根最高(低)作為停損點

//https://www.ptt.cc/bbs/Trading/M.1538318192.A.FBC.html
//binance create order

//Future work
//next strategy
//multiple interval(?)
//Sharpe ratio //http://www.ifuun.com/a2018082215739276/
//https://github.com/markcheno/go-talib/blob/master/talib.go

func main() {
	service.InitStrategyService()

	if core.Config.Trading.Mode == "indicator" {
		service.IndicatorMode()
	} else {
		service.BacktestingMode()
	}
}
