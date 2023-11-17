// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/strategy"
)

//todo
//test slack notify

func main() {
	core.FutureClientInit()
	core.GetHistoryKline()
	//send to channel
	go core.InitWsTickService()
	go strategy.InitStrategyService()

	select {}
}
