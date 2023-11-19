// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/strategy"
)

func main() {
	core.InitSlog()

	strategy.InitStrategyService()
	core.LoadConfigs()
	core.FutureClientInit()
	core.GetHistoryKline()

	go core.InitWsTickService()

	select {}
}
