// strategyService
package strategy

import (
	"log"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	strategySlice []strategyObj
	kLineSlice    []futures.WsKline
)

type strategyObj struct {
	execute func(chan bool)
	notify  chan bool
}

func InitStrategyService(chanNewKline chan futures.WsKline) {
	log.Println("initAnalyzeService() Start")
	initCustomStrategies()

	for {
		newKline := <-chanNewKline
		kLineSlice = append(kLineSlice, newKline)

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
