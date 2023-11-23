// strategyService
package service

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/strategy"
	"fmt"
)

var (
	strategySlice []*core.StrategyBO
)

func InitStrategyService() {
	initCustomStrategies()

	go waitPriceFeeding()

	Logger.Info("[InitStrategyService] Initialized")
}

func initCustomStrategies() {
	strategiesSlice := []func(*core.Kline, *core.StrategyBO){strategy.DoubleTopBottom}

	for idx, strategyFunc := range strategiesSlice {
		strategySlice = append(strategySlice, core.ConstructStrategy(fmt.Sprintf("%d", idx), strategyFunc))
	}
}

func waitPriceFeeding() {
	for {
		nextKline := <-NotifyNewKline

		for _, _strategy := range strategySlice {
			_strategy.GetChanNextKline() <- *nextKline
		}

		for _, _strategy := range strategySlice {
			<-_strategy.GetChanDoneAction()
		}
	}
}
