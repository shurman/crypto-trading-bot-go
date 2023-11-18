// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/strategy"
)

//todo
//test slack notify

func main() {
	strategy.InitStrategyService()

	core.LoadConfigs()
	core.FutureClientInit()
	core.GetHistoryKline()

	go core.InitWsTickService()

	select {}
}
