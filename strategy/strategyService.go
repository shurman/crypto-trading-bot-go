// strategyService
package strategy

import (
	"crypto-trading-bot-go/core"
	"fmt"
	"time"
)

var (
	strategySlice []strategyObj
)

type strategyObj struct {
	name      string
	execute   func(strategyObj)
	nextKline chan core.Kline
}

func InitStrategyService() {
	initCustomStrategies()

	go waitPriceFeeding()

	core.Logger.Info("[InitStrategyService] Initialized")
}

func initCustomStrategies() {
	strategiesSlice := []func(strategyObj){doubleTopBottom}

	for idx, strategyFunc := range strategiesSlice {
		strategySlice = append(strategySlice, constructStrategy(fmt.Sprintf("%d", idx), strategyFunc))
	}
}

func constructStrategy(_name string, _execute func(strategyObj)) strategyObj {
	newStrategy := strategyObj{
		name:      _name,
		execute:   _execute,
		nextKline: make(chan core.Kline),
	}
	go _execute(newStrategy)

	return newStrategy
}

func waitPriceFeeding() {
	for {
		nextKline := <-core.NotifyNewKline

		for _, _strategy := range strategySlice {
			_strategy.nextKline <- *nextKline
		}
	}
}

func (obj *strategyObj) createOrder(
	kline *core.Kline,
	id string,
	dir OrderDirection,
	entry float64,
	stopProfit float64,
	stopLoss float64,
) {
	core.Logger.Info(fmt.Sprintf("[%s][%s][%s][%s] Entry:%f P:%f L:%f",
		time.Unix(kline.StartTime/1000, 0),
		obj.name,
		id,
		dir.toString(),
		entry,
		stopProfit,
		stopLoss))
}

func (obj *strategyObj) cancelOrder(kline *core.Kline, id string) {
	core.Logger.Info(fmt.Sprintf("[%s][%s][%s] cancelled",
		time.Unix(kline.StartTime/1000, 0),
		obj.name,
		id))
}
