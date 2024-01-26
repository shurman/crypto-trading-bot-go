// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//check chart
//https://www.ptt.cc/bbs/Trading/M.1538318192.A.FBC.html
//binance/WOOX create order
//Sharpe ratio //http://www.ifuun.com/a2018082215739276/

//Future work
//backtesting with bitfinex
//multiple interval
//next strategy (Open when EMA alignment trend exists) (ATR as Stop loss) //https://themarketmemo.com/mytradingview/

//backup
//交易不可能三角 //勝率//盈虧比//時間 //https://www.fat-tail.com/v1/selector/article/59e61ccfe1fbca000e853c20.html
//https://github.com/markcheno/go-talib/blob/master/talib.go

func main() {
	service.InitStrategyService()

	if core.Config.Trading.Mode == "indicator" {
		service.IndicatorMode()
	} else {
		service.BacktestingMode()
	}
}
