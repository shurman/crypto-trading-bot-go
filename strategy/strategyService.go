// strategyService
package strategy

import (
	"crypto-trading-bot-go/core"
)

var (
	strategySlice []strategyObj
)

type strategyObj struct {
	execute func(chan bool)
	notify  chan bool
}

func InitStrategyService() {
	initCustomStrategies()

	go waitPriceFeeding()

	core.Logger.Info("[InitStrategyService] Initialized")
}

func initCustomStrategies() {
	strategiesSlice := []func(chan bool){doubleTopBottom}

	for _, strategyFunc := range strategiesSlice {
		strategySlice = append(strategySlice, constructStrategy(strategyFunc, make(chan bool)))
	}
}

func constructStrategy(_execute func(chan bool), _notify chan bool) strategyObj {
	newStrategy := strategyObj{
		execute: _execute,
		notify:  _notify,
	}
	go _execute(_notify)

	return newStrategy
}

func waitPriceFeeding() {
	for {
		<-core.NotifyNewKline

		for _, _strategy := range strategySlice {
			_strategy.notify <- true
		}
	}
}
