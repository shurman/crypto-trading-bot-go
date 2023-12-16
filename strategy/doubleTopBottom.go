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

	mapOrderId = make(map[string]*int)

	lastKline1 = make(map[string]*core.Kline)
	lastKline2 = make(map[string]*core.Kline)
	lastKline3 = make(map[string]*core.Kline)
	lastKline4 = make(map[string]*core.Kline)

	phase2 = true
)

func init() {
	service.RegisterStrategyFunc(DoubleTopBottom, "DTB")
}

func DoubleTopBottom(nextKline *core.Kline, bo *core.StrategyBO) {
	symbol := bo.GetSymbol()

	if nextKline.High == nextKline.Low {
		return
	}

	lastKline4[symbol] = lastKline3[symbol]
	lastKline3[symbol] = lastKline2[symbol]
	lastKline2[symbol] = lastKline1[symbol]
	lastKline1[symbol] = nextKline

	if lastKline3[symbol] == nil {
		paramInit(symbol)
	}
	if lastKline4[symbol] == nil {
		return
	}

	k1 := lastKline1[symbol]
	k2 := lastKline2[symbol]
	k3 := lastKline3[symbol]
	k4 := lastKline4[symbol]

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

		} else if k1.Low > *localLow+(*localHigh-*localLow)*0.15 {
			*state = 2
		}

	} else if *state == 2 {
		if k1.Low < *localLow {
			*state = 3

		} else if k1.High > *localHigh+(*localHigh-*localLow)*0.4 {
			paramReset(symbol)
		}

	} else if *state == 3 {
		//if k1.High > k2.High {
		createLongOrder(bo, k1, false)
		*state = 4
		//}

	} else if *state == 4 {
		if service.GetOrderStatus(bo, genOrderId(symbol, false)) == core.ORDER_ENTRY {
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
			*state = 6
		} else if k1.High > *localHigh+(*localHigh-*localLow)*0.4 {
			paramReset(symbol)
		}

	} else if *state == 6 {
		//if k1.High > k2.High {
		createLongOrder(bo, k1, true)
		*state = 7
		//}
		/*else if k1.Low < *localLow-(*localHigh-*localLow)*0.4 {
			paramReset(symbol)
		}*/

	} else if *state == 7 {
		if service.GetOrderStatus(bo, genOrderId(symbol, true)) == core.ORDER_ENTRY {
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

		} else if k1.High < *localHigh-(*localHigh-*localLow)*0.15 {
			*state = -2
		}

	} else if *state == -2 {
		if k1.High > *localHigh {
			*state = -3

		} else if k1.Low < *localLow-(*localHigh-*localLow)*0.4 {
			paramReset(symbol)
		}

	} else if *state == -3 {
		//if k1.Low < k2.Low {
		createShortOrder(bo, k1, false)
		*state = -4
		//}

	} else if *state == -4 {
		if service.GetOrderStatus(bo, genOrderId(symbol, false)) == core.ORDER_ENTRY {
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
			*state = -6
		} else if k1.Low < *localLow-(*localHigh-*localLow)*0.4 {
			paramReset(symbol)
		}

	} else if *state == -6 {
		//if k1.Low < k2.Low {
		createShortOrder(bo, k1, true)
		*state = -7

		//}
		/*else if k1.High > *localHigh+(*localHigh-*localLow)*0.4 {
			paramReset(symbol)
		}*/

	} else if *state == -7 {
		if service.GetOrderStatus(bo, genOrderId(symbol, true)) == core.ORDER_ENTRY {
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

func getQuantity(symbol string) float64 {
	if core.Config.Trading.EnableAccumulated {
		return service.CurrentFund[symbol] * core.Config.Trading.SingleRiskRatio / (lastKline1[symbol].High - lastKline1[symbol].Low)
	} else {
		return core.Config.Trading.InitialFund * core.Config.Trading.SingleRiskRatio / (lastKline1[symbol].High - lastKline1[symbol].Low)
	}
}

func createLongOrder(bo *core.StrategyBO, kline *core.Kline, isPhase2 bool) {
	service.CreateOrder(
		bo,
		genOrderId(bo.GetSymbol(), isPhase2),
		core.ORDER_LONG,
		getQuantity(bo.GetSymbol()),
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
		getQuantity(bo.GetSymbol()),
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
}
