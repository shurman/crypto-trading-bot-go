// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/strategy"
)

func main() {
	go core.InitWsTickService()
	go strategy.InitStrategyService()

	select {}
}
