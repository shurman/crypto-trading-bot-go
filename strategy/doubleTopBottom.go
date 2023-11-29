// doubleTopBottom
package strategy

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	"fmt"
	// "log/slog"
)

var (
	state     = 0
	borderLow = 0.0
	localHigh = 0.0
	localLow  = 9999999.99

	singleLoss = 30.0

	orderId = 1

	k1 *core.Kline
	k2 *core.Kline
	k3 *core.Kline
)

func init() {
	service.RegisterStrategyFunc(DoubleTopBottom, "DoubleTopBottom")
}

func DoubleTopBottom(nextKline *core.Kline, bo *core.StrategyBO) {
	k3 = k2
	k2 = k1
	k1 = nextKline

	if k3 == nil {
		return
	}

	//slog.Info(fmt.Sprintf("[doubleTopBottom] state=%d %s", state, k1.ToString()))

	if state == 0 {
		if k1.High > k2.High && k2.High > k3.High && k1.High/k3.High > 1.01 {

			state = 1
			localHigh = k1.High
			borderLow = k3.Low
		}

	} else if state == 1 {
		if k1.High > localHigh {
			localHigh = k1.High
			localLow = 9999999.99

		} else if k1.Low < localLow {
			localLow = k1.Low

		} else if k1.Low > localLow+(localHigh-localLow)*0.15 {
			state = 2
		}

	} else if state == 2 {
		if k1.Low < localLow {
			state = 3

		} else if k1.Low < localLow-(localHigh-localLow)*0.4 {
			reset()

		} else if k1.High > localHigh+(localHigh-localLow)*0.4 {
			reset()
		}

	} else if state == 3 {
		if k1.High > k2.High {
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				singleLoss/(k1.High-k1.Low),
				k1.High,
				k1.High+(k1.High-k1.Low),
				k1.Low,
				k1.IsNew,
			)
			//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
			state = 4
		}

	} else if state == 4 {
		if service.GetOrderStatus(bo, genOrderId()) == core.ORDER_ENTRY { // means filled
			//state = 5
			//slog.Info("[Filled]")
			orderId++
			reset()
		} else {
			//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				singleLoss/(k1.High-k1.Low),
				k1.High,
				k1.High+(k1.High-k1.Low),
				k1.Low,
				k1.IsNew,
			)
		}

	} else if state == 5 { //phase 2
		if k1.High < k2.High {
			state = 6
		}

	} else if state == 6 {
		if k1.High > k2.High {
			//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				singleLoss/(k1.High-k1.Low),
				k1.High,
				k1.High+(k1.High-k1.Low),
				k1.Low,
				k1.IsNew,
			)
			state = 7
		}
	} else if state == 7 {
		if k1.High > k2.High { // means filled
			reset()
		} else {
			//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
			service.CreateOrder(
				bo,
				genOrderId(),
				core.ORDER_LONG,
				singleLoss/(k1.High-k1.Low),
				k1.High,
				k1.High+(k1.High-k1.Low),
				k1.Low,
				k1.IsNew,
			)
		}
	}

	//slog.Info(fmt.Sprintf("================= state=%d  localHigh=%f localLow=%f", state, localHigh, localLow))
}

func genOrderId() string {
	return fmt.Sprintf("%04d", orderId)
}

func reset() {
	state = 0
	localHigh = 0
	localLow = 9999999.99
}
