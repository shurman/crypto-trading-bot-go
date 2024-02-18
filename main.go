// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//next indicator//based on OI,LSR

//Future work
//binance/WOOX create order
//Sharpe ratio //http://www.ifuun.com/a2018082215739276/

//backup
//https://www.ptt.cc/bbs/Trading/M.1538318192.A.FBC.html
//交易不可能三角 //勝率//盈虧比//時間 //https://www.fat-tail.com/v1/selector/article/59e61ccfe1fbca000e853c20.html
//https://github.com/markcheno/go-talib/blob/master/talib.go

func main() {
	service.InitStrategyService()

	if core.Config.Trading.Mode == "indicator" {
		service.IndicatorMode()
	} else if core.Config.Trading.Mode == "backtesting" {
		service.BacktestingMode()
	}
}
