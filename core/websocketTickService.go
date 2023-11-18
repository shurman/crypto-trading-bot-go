// websocketTickService
package core

import (
	//"fmt"
	"log/slog"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	NotifyNewKline = make(chan bool)
	lastTick       *futures.WsKline
	currentTick    *futures.WsKline
)

func InitWsTickService() {
	slog.Info("InitWsTickService Start")

	wsKlineHandler := newKlineHandler
	errHandler := func(err error) {
		slog.Error(err.Error())
	}

	doneC, _, err := futures.WsKlineServe(tSymbol, tInterval, wsKlineHandler, errHandler) //WsCombinedKlineServe for multiple
	if err != nil {
		slog.Error(err.Error())
		return
	}
	<-doneC
}

func newKlineHandler(event *futures.WsKlineEvent) {
	//slog.Info(fmt.Sprintf("%+v", event.Kline))

	lastTick = currentTick
	currentTick = &event.Kline

	if lastTick == nil {
		return
	}
	if currentTick.StartTime != lastTick.StartTime {
		newKline := fConvertToKline(lastTick)

		recordNewKline(&newKline)
	}
}

func parseFloat(str string) float64 {
	result, _ := strconv.ParseFloat(str, 10)

	return result
}
