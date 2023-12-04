// klineService
package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
)

var (
// klineSlice []*core.Kline
// klineLen   int = 0
)

// func GetKlineSliceLen() int {
// 	return klineLen
// }

func recordNewKline(symbol string, newKline *core.Kline) {
	CheckOrderFilled(symbol, newKline)

	//klineSlice = append(klineSlice, newKline)
	//klineLen = len(klineSlice)

	Logger.Debug("<- " + fmt.Sprintf("%+v", newKline))
	NotifyNewKline[symbol] <- newKline
	<-NotifyAnalyzeDone[symbol]
}

// func GetLastKline(nth int) *Kline {
// 	if klineLen-nth < 0 {
// 		panic("Index out of range")
// 	}
// 	return klineSlice[klineLen-nth]
// }
