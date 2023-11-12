// crypto-trading-bot-go project main.go
package main

import (
	"github.com/adshao/go-binance/v2/futures"

	"crypto-trading-bot-go/strategy"
)

var (
	newKline = make(chan futures.WsKline)
)

func main() {
	go initWsTickService(newKline)
	go strategy.InitStrategyService(newKline)

	select {}
}
