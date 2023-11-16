// websocketTickService
package core

import (
	//"fmt"
	"log/slog"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	lastTick    futures.WsKline
	currentTick futures.WsKline
)

func InitWsTickService() {
	slog.Info("InitWsTickService() Start")

	wsKlineHandler := func(event *futures.WsKlineEvent) {
		//slog.Info(fmt.Sprintf("%+v", event.Kline))

		sendIfKlineClosed(event.Kline)
	}
	errHandler := func(err error) {
		slog.Error(err.Error())
	}

	doneC, _, err := futures.WsKlineServe("ETHUSDT", "1m", wsKlineHandler, errHandler) //WsCombinedKlineServe for multiple
	if err != nil {
		slog.Error(err.Error())
		return
	}
	<-doneC
}

func sendIfKlineClosed(newTick futures.WsKline) {
	lastTick = currentTick
	currentTick = newTick

	if currentTick.StartTime != lastTick.StartTime {

		newKline := convertToKline(lastTick)

		appendKlineSlice(newKline)
		NotifyNewKline <- true
	}
}

func convertToKline(tick futures.WsKline) Kline {
	return Kline{
		Open:  parseFloat(tick.Open),
		High:  parseFloat(tick.High),
		Low:   parseFloat(tick.Low),
		Close: parseFloat(tick.Close),
	}
}

func parseFloat(str string) float64 {
	result, _ := strconv.ParseFloat(str, 10)

	return result
}
