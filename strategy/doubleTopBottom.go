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

	orderId = 1

	lastKline1 = make(map[string]*core.Kline)
	lastKline2 = make(map[string]*core.Kline)
	lastKline3 = make(map[string]*core.Kline)

	phase2 = true
)

func init() {
	service.RegisterStrategyFunc(DoubleTopBottom, "DoubleTopBottom")
}

func DoubleTopBottom(nextKline *core.Kline, bo *core.StrategyBO) {
	lastKline3[bo.GetSymbol()] = lastKline2[bo.GetSymbol()]
	lastKline2[bo.GetSymbol()] = lastKline1[bo.GetSymbol()]
	lastKline1[bo.GetSymbol()] = nextKline

	if lastKline2[bo.GetSymbol()] == nil {
		paramInit(bo.GetSymbol())
	}
	if lastKline3[bo.GetSymbol()] == nil {
		return
	}

	k1 := lastKline1[bo.GetSymbol()]
	k2 := lastKline2[bo.GetSymbol()]
	k3 := lastKline3[bo.GetSymbol()]

	state := mapState[bo.GetSymbol()]
	borderLow := mapBorderLow[bo.GetSymbol()]
	borderHigh := mapBorderHigh[bo.GetSymbol()]
	localHigh := mapLocalHigh[bo.GetSymbol()]
	localLow := mapLocalLow[bo.GetSymbol()]

	//slog.Info(fmt.Sprintf("[doubleTopBottom] state=%d %s", state, k1.ToString()))

	if *state == 0 {
		if k1.High > k2.High && k2.High > k3.High && k1.High/k3.High > 1.01 {

			*state = 1
			*localHigh = k1.High
			*borderLow = k3.Low
		} else if k1.Low < k2.Low && k2.Low < k3.Low && k1.Low/k3.Low < 0.99 {

			*state = -1
			*localLow = k1.Low
			*borderHigh = k3.High
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

		} else if k1.Low < *localLow-(*localHigh-*localLow)*0.4 {
			paramReset(bo.GetSymbol())

		} else if k1.High > *localHigh+(*localHigh-*localLow)*0.4 {
			paramReset(bo.GetSymbol())
		}

	} else if *state == 3 {
		if k1.High > k2.High {
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				getQuantity(bo.GetSymbol()),
				k1.High,
				k1.High+(k1.High-k1.Low)*core.Config.Trading.ProfitLossRatio,
				k1.Low,
				k1.IsNew,
			)
			//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
			*state = 4
		}

	} else if *state == 4 {
		if service.GetOrderStatus(bo, genOrderId()) == core.ORDER_ENTRY { // means filled
			//state = 5
			//slog.Info("[Filled]")
			orderId++
			paramReset(bo.GetSymbol())
		} else {
			//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				getQuantity(bo.GetSymbol()),
				k1.High,
				k1.High+(k1.High-k1.Low)*core.Config.Trading.ProfitLossRatio,
				k1.Low,
				k1.IsNew,
			)
		}

	} else if *state == 5 { //phase 2
		if k1.High < k2.High {
			*state = 6
		}

	} else if *state == 6 {
		if k1.High > k2.High {
			//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				getQuantity(bo.GetSymbol()),
				k1.High,
				k1.High+(k1.High-k1.Low)*core.Config.Trading.ProfitLossRatio,
				k1.Low,
				k1.IsNew,
			)
			*state = 7
		}
	} else if *state == 7 {
		if k1.High > k2.High { // means filled
			paramReset(bo.GetSymbol())
		} else {
			//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				getQuantity(bo.GetSymbol()),
				k1.High,
				k1.High+(k1.High-k1.Low)*core.Config.Trading.ProfitLossRatio,
				k1.Low,
				k1.IsNew,
			)
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

		} else if k1.High > *localHigh+(*localHigh-*localLow)*0.4 {
			paramReset(bo.GetSymbol())

		} else if k1.Low < *localLow-(*localHigh-*localLow)*0.4 {
			paramReset(bo.GetSymbol())
		}

	} else if *state == -3 {
		if k1.Low < k2.Low {
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_SHORT,
				getQuantity(bo.GetSymbol()),
				k1.Low,
				k1.Low-(k1.High-k1.Low)*core.Config.Trading.ProfitLossRatio,
				k1.High,
				k1.IsNew,
			)
			*state = -4
		}

	} else if *state == -4 {
		if service.GetOrderStatus(bo, genOrderId()) == core.ORDER_ENTRY { // means filled
			//state = -5
			orderId++
			paramReset(bo.GetSymbol())
		} else {
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				getQuantity(bo.GetSymbol()),
				k1.Low,
				k1.Low-(k1.High-k1.Low)*core.Config.Trading.ProfitLossRatio,
				k1.High,
				k1.IsNew,
			)
		}

	}

	//slog.Info(fmt.Sprintf("================= state=%d  CurrentFund=%f localHigh=%f localLow=%f", state, service.CurrentFund, localHigh, localLow))
}

func genOrderId() string {
	return fmt.Sprintf("%04d", orderId)
}

func getQuantity(symbol string) float64 {
	if core.Config.Trading.EnableAccumulated {
		return service.CurrentFund[symbol] * core.Config.Trading.SingleRiskRatio / (lastKline1[symbol].High - lastKline1[symbol].Low)
	} else {
		return core.Config.Trading.InitialFund * core.Config.Trading.SingleRiskRatio / (lastKline1[symbol].High - lastKline1[symbol].Low)
	}
}

func paramInit(symbol string) {
	mapState[symbol] = new(int)
	mapBorderLow[symbol] = new(float64)
	mapBorderHigh[symbol] = new(float64)
	mapLocalHigh[symbol] = new(float64)
	mapLocalLow[symbol] = new(float64)

	paramReset(symbol)
}

func paramReset(symbol string) {
	*mapState[symbol] = 0
	*mapBorderLow[symbol] = 0
	*mapBorderHigh[symbol] = 0
	*mapLocalHigh[symbol] = 0
	*mapLocalLow[symbol] = 9999999.99
}
