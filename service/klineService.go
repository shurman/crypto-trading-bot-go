// klineService
package service

import (
	"crypto-trading-bot-go/core"
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	klineSlice []*core.Kline
	klineLen   int = 0

	NotifyNewKline = make(chan *core.Kline)
)

// func GetKlineSliceLen() int {
// 	return klineLen
// }

func recordNewKline(newKline *core.Kline) {
	klineSlice = append(klineSlice, newKline)
	klineLen = len(klineSlice)

	Logger.Debug("<- " + fmt.Sprintf("%+v", newKline))
	NotifyNewKline <- newKline
}

// func GetLastKline(nth int) *Kline {
// 	if klineLen-nth < 0 {
// 		panic("Index out of range")
// 	}
// 	return klineSlice[klineLen-nth]
// }

func fConvertToKline(tick *futures.WsKline) core.Kline {
	return core.Kline{
		StartTime: tick.StartTime,
		Open:      parseFloat(tick.Open),
		High:      parseFloat(tick.High),
		Low:       parseFloat(tick.Low),
		Close:     parseFloat(tick.Close),
		IsNew:     true,
	}
}

func kConvertKline(tick *futures.Kline) core.Kline {
	return core.Kline{
		StartTime: tick.OpenTime,
		Open:      parseFloat(tick.Open),
		High:      parseFloat(tick.High),
		Low:       parseFloat(tick.Low),
		Close:     parseFloat(tick.Close),
		IsNew:     false,
	}
}
