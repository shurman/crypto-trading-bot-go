// doubleTopBottom
package strategy

import (
	"crypto-trading-bot-go/core"
	"fmt"
	"log/slog"
	"time"
)

var (
	state     = 0
	borderLow = 0.0
	localHigh = 0.0
	localLow  = 9999999.99

	k1 *core.Kline
	k2 *core.Kline
	k3 *core.Kline
)

func DoubleTopBottom(obj core.StrategyObj) {

	for {
		nextKline := <-obj.NextKline

		k3 = k2
		k2 = k1
		k1 = &nextKline

		if k3 == nil {
			continue
		}

		slog.Info(fmt.Sprintf("[doubleTopBottom] state=%d %s %s", state, time.Unix(k1.StartTime/1000, 0), k1.ToString()))

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
				obj.CreateOrder(
					k1,
					"L1",
					core.LONG,
					1,
					k1.High,
					k1.High+(k1.High-k1.Low),
					k1.Low,
				)
				//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
				state = 4
			}

		} else if state == 4 {
			if k1.High > k2.High { // means filled
				//state = 5
				slog.Info("[Filled]")
				reset()
			} else {
				//cancel order
				obj.CancelOrder(k1, "L1")
				//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
				obj.CreateOrder(
					k1,
					"L1",
					core.LONG,
					1,
					k1.High,
					k1.High+(k1.High-k1.Low),
					k1.Low,
				)
			}

		} else if state == 5 { //phase 2
			if k1.High < k2.High {
				state = 6
			}

		} else if state == 6 {
			if k1.High > k2.High {
				//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
				obj.CreateOrder(
					k1,
					"L2",
					core.LONG,
					1,
					k1.High,
					k1.High+(k1.High-k1.Low),
					k1.Low,
				)
				state = 7
			}
		} else if state == 7 {
			if k1.High > k2.High { // means filled
				reset()
			} else {
				//cancel order
				obj.CancelOrder(k1, "L2")
				//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
				obj.CreateOrder(
					k1,
					"L2",
					core.LONG,
					1,
					k1.High,
					k1.High+(k1.High-k1.Low),
					k1.Low,
				)
			}
		}

		slog.Info(fmt.Sprintf("================= state=%d  localHigh=%f localLow=%f", state, localHigh, localLow))
	}
}

func reset() {
	state = 0
	localHigh = 0
	localLow = 9999999.99
}
