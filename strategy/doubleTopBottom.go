// doubleTopBottom
package strategy

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	"fmt"
	"math"

	"github.com/cinar/indicator"
)

var (
	mapState      = make(map[string]*int)
	mapBorderLow  = make(map[string]*float64)
	mapBorderHigh = make(map[string]*float64)
	mapLocalHigh  = make(map[string]*float64)
	mapLocalLow   = make(map[string]*float64)
	mapOrderId    = make(map[string]*int)

	mapBBCounter = make(map[string]*int)

	phase2           = true
	exitRatio        = 0.4
	borderValidRatio = 0.15

	maxFeeRatio = 0.1

	klinesBreakBase  = 13
	klinesBreakLimit = 7
	klinesLimit      = 20 + klinesBreakBase + 1

	klinesStopLossRefLimit = 5
)

func init() {
	service.RegisterStrategyFunc(DoubleTopBottom, "DTB")
}

func DoubleTopBottom(nextKline *core.Kline, bo *core.StrategyBO) {
	symbol := bo.GetSymbol()

	if service.GetKlinesLen(symbol) == 1 {
		paramInit(symbol)
	}
	if service.GetKlinesLen(symbol) < klinesLimit {
		return
	}

	klines := service.GetRecentKlines(klinesLimit, symbol)
	closedKlines := service.GetClosedPrices(klines)

	_, upperBand, lowerBand := indicator.BollingerBands(closedKlines)

	k1 := klines[klinesLimit-1]
	k2 := klines[klinesLimit-2]
	//k3 := klines[klinesLimit-3]
	//k4 := klines[klinesLimit-4]
	kTail := klines[klinesLimit-klinesBreakBase]

	state := mapState[symbol]
	borderLow := mapBorderLow[symbol]
	borderHigh := mapBorderHigh[symbol]
	localHigh := mapLocalHigh[symbol]
	localLow := mapLocalLow[symbol]

	if *state != 0 {
		service.Logger.Debug(fmt.Sprintf("%s\ts=%d\tborderLow=%f\tborderHigh=%f\tlocalHigh=%f\tlocalLow=%f", k1.StartTime, *state, *borderLow, *borderHigh, *localHigh, *localLow))
	}

	if *state == 0 {
		//========3========
		countBreakUp, countBreakDown := countBreakingKlines(closedKlines[klinesLimit-klinesBreakBase:], upperBand[klinesLimit-klinesBreakBase:], lowerBand[klinesLimit-klinesBreakBase:])
		if countBreakUp >= klinesBreakLimit {
			*state = 1
			*localHigh = k1.High
			*borderLow = kTail.Low
		} else if countBreakDown >= klinesBreakLimit {
			*state = -1
			*localLow = k1.Low
			*borderHigh = kTail.High
		}

		//========2========
		// if k1.High > upperBand[klinesLimit-1] {
		// 	if *bbCounter < 0 {
		// 		*bbCounter = 1
		// 	} else {
		// 		*bbCounter++
		// 	}

		// 	if *bbCounter == klinesBreakLimit {
		// 		*state = 1
		// 		*localHigh = k1.High
		// 		*borderLow = kTail.Low
		// 	}
		// } else if k1.Low < lowerBand[klinesLimit-1] {
		// 	if *bbCounter > 0 {
		// 		*bbCounter = -1
		// 	} else {
		// 		*bbCounter--
		// 	}

		// 	if *bbCounter == -klinesBreakLimit {
		// 		*state = -1
		// 		*localLow = k1.Low
		// 		*borderHigh = kTail.High
		// 	}
		// } else {
		// 	*bbCounter = 0
		// }

		//========1========
		// if k1.High > k2.High && k2.High > k3.High && k3.High > k4.High && k1.High/k4.High > 1.02 {
		// 	*state = 1
		// 	*localHigh = k1.High
		// 	*borderLow = k4.Low

		// } else if k1.Low < k2.Low && k2.Low < k3.Low && k3.Low < k4.Low && k1.Low/k4.Low < 0.98 {
		// 	*state = -1
		// 	*localLow = k1.Low
		// 	*borderHigh = k4.High
		// }

	} else if *state == 1 {
		if k1.High > *localHigh {
			*localHigh = k1.High
			*localLow = 9999999.99

		} else if k1.Low < *localLow {
			*localLow = k1.Low

		} else if k1.Low < *borderLow {
			paramReset(symbol)

		} else if k1.Low > *localLow+(*localHigh-*localLow)*borderValidRatio {
			*state = 2
		}

	} else if *state == 2 {
		if k1.Low < *localLow {
			createLongOrder(bo, k1, false)
			//*state = 3
			*state = 4

		} else if k1.High > *localHigh+(*localHigh-*localLow)*exitRatio {
			*state = 1
			*localHigh = k1.High
			*localLow = 9999999.99
			//paramReset(symbol)
		} else if k1.High > upperBand[len(upperBand)-1] {
			paramReset(symbol)
		}

	} else if *state == 3 {
		//createLongOrder(bo, k1, false)
		//*state = 4

	} else if *state == 4 {
		if service.GetOrder(bo, genOrderId(symbol, false)).IsExistedAndFilled() {
			if phase2 {
				*state = 5
			} else {
				paramReset(symbol)
			}
			incrOrderId(symbol)

		} else if k1.Low < *borderLow {
			service.CancelOrder(bo, genOrderId(symbol, false), k1.IsNew)
			paramReset(symbol)
		} else {
			createLongOrder(bo, k1, false)
		}

	} else if *state == 5 { //phase 2
		if k1.High < k2.High {
			createLongOrder(bo, k1, true)
			*state = 7
			//*state = 6
		} else if k1.High > *localHigh+(*localHigh-*localLow)*exitRatio {
			paramReset(symbol)
		}

	} else if *state == 6 {
		// createLongOrder(bo, k1, true)
		// *state = 7

	} else if *state == 7 {
		if service.GetOrder(bo, genOrderId(symbol, true)).IsExistedAndFilled() {
			incrOrderId(symbol)
			paramReset(symbol)
		} else {
			createLongOrder(bo, k1, true)
		}

	} else if *state == -1 {
		if k1.Low < *localLow {
			*localHigh = 0
			*localLow = k1.Low

		} else if k1.High > *localHigh {
			*localHigh = k1.High

		} else if k1.High > *borderHigh {
			paramReset(symbol)

		} else if k1.High < *localHigh-(*localHigh-*localLow)*borderValidRatio {
			*state = -2
		}

	} else if *state == -2 {
		if k1.High > *localHigh {
			createShortOrder(bo, k1, false)
			*state = -4
			// *state = -3

		} else if k1.Low < *localLow-(*localHigh-*localLow)*exitRatio {
			//paramReset(symbol)
			*state = -1
			*localHigh = 0
			*localLow = k1.Low
		} else if k1.Low < lowerBand[len(lowerBand)-1] {
			paramReset(symbol)
		}

	} else if *state == -3 {
		// createShortOrder(bo, k1, false)
		// *state = -4

	} else if *state == -4 {
		if service.GetOrder(bo, genOrderId(symbol, false)).IsExistedAndFilled() {
			if phase2 {
				*state = -5
			} else {
				paramReset(symbol)
			}
			incrOrderId(symbol)

		} else if k1.High > *borderHigh {
			service.CancelOrder(bo, genOrderId(symbol, false), k1.IsNew)
			paramReset(symbol)
		} else {
			createShortOrder(bo, k1, false)
		}

	} else if *state == -5 { //phase 2
		if k1.Low > k2.Low {
			createShortOrder(bo, k1, true)
			*state = -7
			// *state = -6
		} else if k1.Low < *localLow-(*localHigh-*localLow)*exitRatio {
			paramReset(symbol)
		}

	} else if *state == -6 {
		// createShortOrder(bo, k1, true)
		// *state = -7

	} else if *state == -7 {
		if service.GetOrder(bo, genOrderId(symbol, true)).IsExistedAndFilled() {
			incrOrderId(symbol)
			paramReset(symbol)
		} else {
			createShortOrder(bo, k1, true)
		}
	}
}

func genOrderId(symbol string, isPhase2 bool) string {
	orderId := fmt.Sprintf("%s%05d", symbol, *mapOrderId[symbol])

	if isPhase2 {
		orderId += "_p2"
	}
	return orderId
}

func incrOrderId(symbol string) {
	*mapOrderId[symbol]++
}

func getQuantity(symbol string, entryPrice float64, stopLossPrice float64) float64 {
	return service.GetRiskPerTrade(symbol) / math.Abs(entryPrice-stopLossPrice)
}

func createLongOrder(bo *core.StrategyBO, kline *core.Kline, isPhase2 bool) {
	if kline.High == kline.Low {
		return
	}

	//Prevent fee exceeding limit
	if kline.High-kline.Low < kline.High*(core.Config.Trading.FeeTakerRate/100)/maxFeeRatio {
		return
	}

	_, low := getKlinesHighLow(service.GetRecentKlines(klinesStopLossRefLimit, bo.GetSymbol()))

	service.CreateOrder(
		bo,
		genOrderId(bo.GetSymbol(), isPhase2),
		core.ORDER_LONG,
		getQuantity(bo.GetSymbol(), kline.High, low),
		kline.High,
		kline.High+(kline.High-low)*core.Config.Trading.ProfitLossRatio,
		low,
		kline.IsNew,
	)
}

func createShortOrder(bo *core.StrategyBO, kline *core.Kline, isPhase2 bool) {
	if kline.High == kline.Low {
		return
	}

	if kline.High-kline.Low < kline.Low*(core.Config.Trading.FeeTakerRate/100)/maxFeeRatio {
		return
	}

	high, _ := getKlinesHighLow(service.GetRecentKlines(klinesStopLossRefLimit, bo.GetSymbol()))

	service.CreateOrder(
		bo,
		genOrderId(bo.GetSymbol(), isPhase2),
		core.ORDER_SHORT,
		getQuantity(bo.GetSymbol(), kline.Low, high),
		kline.Low,
		kline.Low-(high-kline.Low)*core.Config.Trading.ProfitLossRatio,
		high,
		kline.IsNew,
	)
}

func paramInit(symbol string) {
	mapState[symbol] = new(int)
	mapBorderLow[symbol] = new(float64)
	mapBorderHigh[symbol] = new(float64)
	mapLocalHigh[symbol] = new(float64)
	mapLocalLow[symbol] = new(float64)
	mapBBCounter[symbol] = new(int)

	mapOrderId[symbol] = new(int)
	*mapOrderId[symbol] = 1

	paramReset(symbol)
}

func paramReset(symbol string) {
	*mapState[symbol] = 0
	*mapBorderLow[symbol] = 0
	*mapBorderHigh[symbol] = 0
	*mapLocalHigh[symbol] = 0
	*mapLocalLow[symbol] = 9999999.99
	*mapBBCounter[symbol] = 0
}

func countBreakingKlines(closedKline []float64, upperBand []float64, lowerBand []float64) (int, int) {
	countBreakUp := 0
	countBreakDown := 0

	for idx, price := range closedKline {
		if price > upperBand[idx] {
			countBreakUp += 1
		} else if price < lowerBand[idx] {
			countBreakDown += 1
		}
	}

	return countBreakUp, countBreakDown
}

func getKlinesHighLow(klines []*core.Kline) (high float64, low float64) {
	high = 0.0
	low = math.MaxFloat64
	for _, k := range klines {
		if k.High > high {
			high = k.High
		}
		if k.Low < low {
			low = k.Low
		}
	}
	return
}
