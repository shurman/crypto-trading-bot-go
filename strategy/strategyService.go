// strategyService
package strategy

import (
	"log"

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
	log.Println("InitStrategyService() Start")
	initCustomStrategies()

	for {
		<-core.NotifyNewKline

		for _, _strategy := range strategySlice {
			_strategy.notify <- true
		}
	}
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
