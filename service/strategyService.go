// strategyService
package service

import (
	"crypto-trading-bot-go/core"
)

var (
	strategySlice []*core.StrategyBO
)

func InitStrategyService() {
	go waitPriceFeeding()

	Logger.Info("[InitStrategyService] Initialized")
}

func RegisterStrategyFunc(f func(*core.Kline, *core.StrategyBO), name string) {
	strategySlice = append(strategySlice, core.ConstructStrategy(name, f))
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

		NotifyKlineDone <- true
	}
}
