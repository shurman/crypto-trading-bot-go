// klineService
package core

import (
	"github.com/adshao/go-binance/v2/futures"
)

var (
	klineSlice []*Kline
	klineLen   int = 0
)

type Kline struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

func GetKlineSliceLen() int {
	return klineLen
}

func recordNewKline(newKline *Kline) {
	klineSlice = append(klineSlice, newKline)
	klineLen = len(klineSlice)

	//slog.Info(fmt.Sprintf("%+v", newKline))
	NotifyNewKline <- true
}

func GetLastKline(nth int) *Kline {
	if klineLen-nth < 0 {
		panic("Index out of range")
	}
	return klineSlice[klineLen-nth]
}

func fConvertToKline(tick *futures.WsKline) Kline {
	return Kline{
		Open:  parseFloat(tick.Open),
		High:  parseFloat(tick.High),
		Low:   parseFloat(tick.Low),
		Close: parseFloat(tick.Close),
	}
}

func kConvertKline(tick *futures.Kline) Kline {
	return Kline{
		Open:  parseFloat(tick.Open),
		High:  parseFloat(tick.High),
		Low:   parseFloat(tick.Low),
		Close: parseFloat(tick.Close),
	}
}
