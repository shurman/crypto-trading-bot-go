// klineService
package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
)

var (
	symbolKlines = make(map[string][]*core.Kline)
)

func init() {
	// for _, symbol := range core.Config.Trading.Symbols {
	// 	symbolKlines[symbol]
	// }
}

func recordNewKline(symbol string, newKline *core.Kline) {
	CheckOrderFilled(symbol, newKline)

	symbolKlines[symbol] = append(symbolKlines[symbol], newKline)

	Logger.Debug("<- " + fmt.Sprintf("%+v", newKline))
	NotifyNewKline[symbol] <- newKline
	<-NotifyAnalyzeDone[symbol]
}

func GetKlinesLen(symbol string) int {
	return len(symbolKlines[symbol])
}

func GetRecentKlines(limit int, symbol string) []*core.Kline {
	if limit > 40 {
		Logger.Warn("Cannot load more than 30 klines")
		return nil
	} else if GetKlinesLen(symbol)-limit < 0 {
		Logger.Warn("Index out of range")
		return nil
	}
	return symbolKlines[symbol][GetKlinesLen(symbol)-limit:]
}

func GetClosedPrices(rawKlines []*core.Kline) []float64 {
	var closedKlines []float64
	for _, kline := range rawKlines {
		closedKlines = append(closedKlines, kline.Close)
	}

	return closedKlines
}
