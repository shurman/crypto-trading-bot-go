// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//回檔必須超過均線
//入場後被震盪掃出場: 找近幾根最高(低)作為停損點 or 支撐線
//backup: state 2 突破時 如果突破BB範圍 則取消(?

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
