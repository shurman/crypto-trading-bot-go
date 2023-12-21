// doubleTopBottom
package strategy

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	"fmt"
	//"log/slog"
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
)

func init() {
	service.RegisterStrategyFunc(DoubleTopBottom, "DTB")

	//indicator.BollingerBands()
}

func DoubleTopBottom(nextKline *core.Kline, bo *core.StrategyBO) {
	symbol := bo.GetSymbol()

	if service.GetKlinesLen(symbol) == 1 {
		paramInit(symbol)
	}
	if service.GetKlinesLen(symbol) < 4 {
		return
	}

	klines := service.GetRecentKline(4, symbol)

	k1 := klines[3]
	k2 := klines[2]
	k3 := klines[1]
	k4 := klines[0]

	state := mapState[symbol]
	borderLow := mapBorderLow[symbol]
	borderHigh := mapBorderHigh[symbol]
	localHigh := mapLocalHigh[symbol]
	localLow := mapLocalLow[symbol]

	//slog.Info(fmt.Sprintf("[doubleTopBottom] state=%d %s", state, k1.ToString()))

	if *state == 0 {
		if k1.High > k2.High && k2.High > k3.High && k3.High > k4.High && k1.High/k4.High > 1.02 {
			*state = 1
			*localHigh = k1.High
			*borderLow = k4.Low

		} else if k1.Low < k2.Low && k2.Low < k3.Low && k3.Low < k4.Low && k1.Low/k4.Low < 0.98 {
			*state = -1
			*localLow = k1.Low
			*borderHigh = k4.High
		}

	} else if *state == 1 {
		if k1.High > *localHigh {
			*localHigh = k1.High
			*localLow = 9999999.99

		} else if k1.Low < *localLow {
			*localLow = k1.Low

		} else if k1.Low > *localLow+(*localHigh-*localLow)*borderValidRatio {
			*state = 2
		}

	} else if *state == 2 {
		if k1.Low < *localLow {
			createLongOrder(bo, k1, false)
			//*state = 3
			*state = 4

		} else if k1.High > *localHigh+(*localHigh-*localLow)*exitRatio {
			paramReset(symbol)
		}

	} else if *state == 3 {
		//createLongOrder(bo, k1, false)
		//*state = 4

	} else if *state == 4 {
		if service.GetOrder(bo, genOrderId(symbol, false)).IsFilled() {
			if phase2 {
				*state = 5
			} else {
				paramReset(symbol)
			}
			incrOrderId(symbol)

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
		if service.GetOrder(bo, genOrderId(symbol, true)).IsFilled() {
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

		} else if k1.High < *localHigh-(*localHigh-*localLow)*borderValidRatio {
			*state = -2
		}

	} else if *state == -2 {
		if k1.High > *localHigh {
			createShortOrder(bo, k1, false)
			*state = -4
			// *state = -3

		} else if k1.Low < *localLow-(*localHigh-*localLow)*exitRatio {
			paramReset(symbol)
		}

	} else if *state == -3 {
		// createShortOrder(bo, k1, false)
		// *state = -4

	} else if *state == -4 {
		if service.GetOrder(bo, genOrderId(symbol, false)).IsFilled() {
			if phase2 {
				*state = -5
			} else {
				paramReset(symbol)
			}
			incrOrderId(symbol)

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
		if service.GetOrder(bo, genOrderId(symbol, true)).IsFilled() {
			incrOrderId(symbol)
			paramReset(symbol)
		} else {
			createShortOrder(bo, k1, true)
		}
	}
	//slog.Info(fmt.Sprintf("================= state=%d  CurrentFund=%f localHigh=%f localLow=%f", state, service.CurrentFund, localHigh, localLow))
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

func getQuantity(symbol string, lastKline *core.Kline) float64 {
	if core.Config.Trading.EnableAccumulated {
		return service.CurrentFund[symbol] * core.Config.Trading.SingleRiskRatio / (lastKline.High - lastKline.Low)
	} else {
		return core.Config.Trading.InitialFund * core.Config.Trading.SingleRiskRatio / (lastKline.High - lastKline.Low)
	}
}

func createLongOrder(bo *core.StrategyBO, kline *core.Kline, isPhase2 bool) {
	service.CreateOrder(
		bo,
		genOrderId(bo.GetSymbol(), isPhase2),
		core.ORDER_LONG,
		getQuantity(bo.GetSymbol(), kline),
		kline.High,
		kline.High+(kline.High-kline.Low)*core.Config.Trading.ProfitLossRatio,
		kline.Low,
		kline.IsNew,
	)
}

func createShortOrder(bo *core.StrategyBO, kline *core.Kline, isPhase2 bool) {
	service.CreateOrder(
		bo,
		genOrderId(bo.GetSymbol(), isPhase2),
		core.ORDER_SHORT,
		getQuantity(bo.GetSymbol(), kline),
		kline.Low,
		kline.Low-(kline.High-kline.Low)*core.Config.Trading.ProfitLossRatio,
		kline.High,
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
