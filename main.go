// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//flag to delay tp,sl for one kline (exit price should not use pre-set price)
//Sharpe ratio //http://www.ifuun.com/a2018082215739276/
//multiple interval

//https://www.ptt.cc/bbs/Trading/M.1538318192.A.FBC.html
//binance/WOOX create order
//backtesting with bitfinex

//Future work
//next strategy
//https://github.com/markcheno/go-talib/blob/master/talib.go

func main() {
	service.InitStrategyService()

	if core.Config.Trading.Mode == "indicator" {
		service.IndicatorMode()
	} else {
		service.BacktestingMode()
	}
}
