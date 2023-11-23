// klineService
package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	//klineSlice []*core.Kline
	//klineLen   int = 0

	NotifyNewKline = make(chan *core.Kline)
)

// func GetKlineSliceLen() int {
// 	return klineLen
// }

func recordNewKline(newKline *core.Kline) {
	core.SetCurrentKline(newKline)
	core.CheckOrderFilled() //TODO

	//klineSlice = append(klineSlice, newKline)
	//klineLen = len(klineSlice)

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
		StartTime: time.Unix(tick.StartTime/1000, 0),
		Open:      parseFloat(tick.Open),
		High:      parseFloat(tick.High),
		Low:       parseFloat(tick.Low),
		Close:     parseFloat(tick.Close),
		CloseTime: time.Unix(tick.EndTime/1000, 0),
		IsNew:     true,
	}
}

func kConvertKline(tick *futures.Kline) core.Kline {
	return core.Kline{
		StartTime: time.Unix(tick.OpenTime/1000, 0),
		Open:      parseFloat(tick.Open),
		High:      parseFloat(tick.High),
		Low:       parseFloat(tick.Low),
		Close:     parseFloat(tick.Close),
		CloseTime: time.Unix(tick.CloseTime/1000, 0),
		IsNew:     false,
	}
}
