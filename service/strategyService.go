// strategyService
package service

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/strategy"
	"fmt"
)

var (
	strategySlice []core.StrategyObj
)

func InitStrategyService() {
	initCustomStrategies()

	go waitPriceFeeding()

	Logger.Info("[InitStrategyService] Initialized")
}

func initCustomStrategies() {
	strategiesSlice := []func(core.StrategyObj){strategy.DoubleTopBottom}

	for idx, strategyFunc := range strategiesSlice {
		strategySlice = append(strategySlice, constructStrategy(fmt.Sprintf("%d", idx), strategyFunc))
	}
}

func constructStrategy(_name string, _execute func(core.StrategyObj)) core.StrategyObj {
	newStrategy := core.StrategyObj{
		Name:      _name,
		Execute:   _execute,
		NextKline: make(chan core.Kline),
	}
	go _execute(newStrategy)

	return newStrategy
}

func waitPriceFeeding() {
	for {
		nextKline := <-NotifyNewKline

		for _, _strategy := range strategySlice {
			_strategy.NextKline <- *nextKline
		}
	}
}
