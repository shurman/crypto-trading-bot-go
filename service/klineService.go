// klineService
package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
)

var (
	//klineSlice []*core.Kline
	//klineLen   int = 0

	NotifyNewKline  = make(chan *core.Kline)
	NotifyKlineDone = make(chan bool)
)

// func GetKlineSliceLen() int {
// 	return klineLen
// }

func recordNewKline(newKline *core.Kline) {
	SetCurrentKline(newKline)
	CheckOrderFilled()

	//klineSlice = append(klineSlice, newKline)
	//klineLen = len(klineSlice)

	Logger.Debug("<- " + fmt.Sprintf("%+v", newKline))
	NotifyNewKline <- newKline
	<-NotifyKlineDone
}

// func GetLastKline(nth int) *Kline {
// 	if klineLen-nth < 0 {
// 		panic("Index out of range")
// 	}
// 	return klineSlice[klineLen-nth]
// }
