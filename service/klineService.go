// klineService
package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
)

var (
	//klineSlice []*core.Kline
	//klineLen   int = 0

	NotifyNewKline  = make(map[string]chan *core.Kline)
	NotifyKlineDone = make(map[string]chan bool)
)

type KlineNotify struct {
	Kline  *core.Kline
	Symbol string
}

// func GetKlineSliceLen() int {
// 	return klineLen
// }

func recordNewKline(symbol string, newKline *core.Kline) {
	SetCurrentKline(symbol, newKline)
	CheckOrderFilled()

	//klineSlice = append(klineSlice, newKline)
	//klineLen = len(klineSlice)

	Logger.Debug("<- " + fmt.Sprintf("%+v", newKline))
	NotifyNewKline[symbol] <- newKline
	<-NotifyKlineDone[symbol]
}

// func GetLastKline(nth int) *Kline {
// 	if klineLen-nth < 0 {
// 		panic("Index out of range")
// 	}
// 	return klineSlice[klineLen-nth]
// }
