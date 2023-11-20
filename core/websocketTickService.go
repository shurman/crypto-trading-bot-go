// websocketTickService
package core

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	NotifyNewKline = make(chan bool)
	lastTick       *futures.WsKline
	currentTick    *futures.WsKline
)

func InitWsTickService() {
	wsKlineHandler := newKlineHandler
	errHandler := func(err error) {
		Logger.Error(err.Error())
	}

	doneC, _, err := futures.WsKlineServe( //WsCombinedKlineServe for multiple
		Config.Trading.Symbol,
		Config.Trading.Interval,
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

	lastTick = currentTick
	currentTick = &event.Kline

	if lastTick == nil {
		return
	}
	if currentTick.StartTime != lastTick.StartTime {
		newKline := fConvertToKline(lastTick)

		Logger.Info(fmt.Sprintf("%+v", newKline))
		recordNewKline(&newKline)
	}
}

func parseFloat(str string) float64 {
	result, _ := strconv.ParseFloat(str, 10)

	return result
}
