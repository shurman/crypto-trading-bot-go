// strategyService
package strategy

import (
	"crypto-trading-bot-go/core"
)

var (
	strategySlice []strategyObj
)

type strategyObj struct {
	execute    func(strategyObj)
	notifyNew  chan bool
	notifyDone chan bool
}

func InitStrategyService() {
	initCustomStrategies()

	go waitPriceFeeding()

	core.Logger.Info("[InitStrategyService] Initialized")
}

func initCustomStrategies() {
	strategiesSlice := []func(strategyObj){doubleTopBottom}

	for _, strategyFunc := range strategiesSlice {
		strategySlice = append(strategySlice, constructStrategy(strategyFunc))
	}
}

func constructStrategy(_execute func(strategyObj)) strategyObj {
	newStrategy := strategyObj{
		execute:    _execute,
		notifyNew:  make(chan bool),
		notifyDone: make(chan bool),
	}
	go _execute(newStrategy)

	return newStrategy
}

func waitPriceFeeding() {
	for {
		<-core.NotifyNewKline

		for _, _strategy := range strategySlice {
			_strategy.notifyNew <- true
			<-_strategy.notifyDone
		}

		core.NotifyDone <- true
	}
}
