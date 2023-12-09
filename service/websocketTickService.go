// websocketTickService
package service

import (
	"crypto-trading-bot-go/core"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	lastTick    = make(map[string]*futures.WsKline)
	currentTick = make(map[string]*futures.WsKline)
)

func InitWsTickService() {
	wsKlineHandler := newKlineHandler
	errHandler := func(err error) {
		Logger.Error(err.Error())
	}

	// doneC, _, err := futures.WsKlineServe(
	// 	core.Config.Trading.Symbol,
	// 	core.Config.Trading.Interval,
	// 	wsKlineHandler,
	// 	errHandler,
	// )
	doneC, _, err := futures.WsCombinedKlineServe(
		genSymbolsMap(),
		wsKlineHandler,
		errHandler,
	)
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	Logger.Info("[InitWsTickService] Initialized")
	<-doneC
}

func newKlineHandler(event *futures.WsKlineEvent) {
	Logger.Debug("ws event: " + fmt.Sprintf("%+v", event.Kline))

	lastTick[event.Symbol] = currentTick[event.Symbol]
	currentTick[event.Symbol] = &event.Kline

	if lastTick[event.Symbol] == nil {
		return
	}
	if currentTick[event.Symbol].StartTime != lastTick[event.Symbol].StartTime {
		newKline := core.ConvertToKlineFromWsKline(lastTick[event.Symbol])

		Logger.Info(fmt.Sprintf("%s %s", event.Symbol, newKline.ToString()))
		recordNewKline(event.Symbol, &newKline)
	}
}

func genSymbolsMap() map[string]string {
	symbolMap := make(map[string]string)
	for _, symbol := range core.Config.Trading.Symbols {
		symbolMap[symbol] = core.Config.Trading.Interval
	}

	return symbolMap
}
