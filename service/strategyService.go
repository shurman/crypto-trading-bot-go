// strategyService
package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
)

var (
	strategiesFuncSlice []func(*core.Kline, *core.StrategyBO)
	strategySlice       []*core.StrategyBO
)

func InitStrategyService() {
	initCustomStrategies()

	go waitPriceFeeding()

	Logger.Info("[InitStrategyService] Initialized")
}

func RegisterStrategyFunc(f func(*core.Kline, *core.StrategyBO)) {
	strategiesFuncSlice = append(strategiesFuncSlice, f)
}

func initCustomStrategies() {
	for idx, strategyFunc := range strategiesFuncSlice {
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

		NotifyKlineDone <- true
	}
}
