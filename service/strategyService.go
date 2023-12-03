// strategyService
package service

import (
	"crypto-trading-bot-go/core"
)

var (
	strategyBaseList []*core.StrategyBase
	strategiesMap    = make(map[string]([]*core.StrategyBO))
)

func InitStrategyService() {
	for _, symbol := range core.Config.Trading.Symbols {
		for _, o := range strategyBaseList {
			strategiesMap[symbol] = append(strategiesMap[symbol], core.InitialStrategy(*o, symbol))
		}

		go waitPriceFeeding(symbol)
	}

	Logger.Info("[InitStrategyService] Initialized")
}

func RegisterStrategyFunc(f func(*core.Kline, *core.StrategyBO), name string) {
	strategyBaseList = append(strategyBaseList, core.ConstructStrategyBase(name, f))
}

func waitPriceFeeding(symbol string) {
	initNotifyChan(symbol)

	for {
		nextKline := <-NotifyNewKline[symbol]

		for _, _strategy := range strategiesMap[symbol] {
			_strategy.GetChanNextKline() <- *nextKline
		}

		for _, _strategy := range strategiesMap[symbol] {
			<-_strategy.GetChanDoneAction()
		}

		NotifyKlineDone[symbol] <- true
	}
}

func initNotifyChan(symbol string) {
	NotifyNewKline[symbol] = make(chan *core.Kline)
	NotifyKlineDone[symbol] = make(chan bool)
}
