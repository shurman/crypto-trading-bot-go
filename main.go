// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/strategy"
)

//TODO
//Download history kline and load
//verify performance

func main() {
	core.LoadConfigs()
	core.InitSlog()

	strategy.InitStrategyService()

	core.FutureClientInit()
	core.LoadHistoryKline()
	go core.InitWsTickService()

	select {}
}
