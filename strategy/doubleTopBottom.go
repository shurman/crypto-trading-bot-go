// doubleTopBottom
package strategy

import (
	"crypto-trading-bot-go/core"
	"fmt"
)

var (
	state     = 0
	borderLow = 0.0
	localHigh = 0.0
	localLow  = 9999999.99
)

func doubleTopBottom(notify chan bool) {

	for {
		<-notify

		if core.GetKlineSliceLen() < 3 {
			continue
		}

		k1 := core.GetLastKline(1)
		k2 := core.GetLastKline(2)
		k3 := core.GetLastKline(3)

		core.Logger.Info(fmt.Sprintf("[doubleTopBottom] ==== state=%d  %+v", state, k1))

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

			} else if k1.Low > (localHigh-localLow)*0.15 {
				state = 2
			}

		} else if state == 2 {
			if k1.Low < localLow {
				state = 3

			} else if k1.Low < borderLow-(localHigh-localLow)*0.4 {
				reset()

			} else if k1.High > localHigh+(localHigh-localLow)*0.4 {
				reset()
			}

		} else if state == 3 {
			if k1.High > k2.High {
				core.Logger.Info("[Place Long] Trigger@" + fmt.Sprintf("%f", k1.High) +
					" P:" + fmt.Sprintf("%f", k1.High+(k1.High-k1.Low)) +
					" L:" + fmt.Sprintf("%f", k1.High))
				//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
				state = 4
			}

		} else if state == 4 {
			if k1.High > k2.High { // means filled
				//state = 5
				core.Logger.Info("[Filled]")
				reset()
			} else {
				//cancel order
				//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
				core.Logger.Info("[Re-Place Long] Trigger@" + fmt.Sprintf("%f", k1.High) +
					" P:" + fmt.Sprintf("%f", k1.High+(k1.High-k1.Low)) +
					" L:" + fmt.Sprintf("%f", k1.High))
			}

		} else if state == 5 { //phase 2
			if k1.High < k2.High {
				state = 6
			}

		} else if state == 6 {
			if k1.High > k2.High {
				//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
				core.Logger.Info("[Place Long] Trigger@" + fmt.Sprintf("%f", k1.High) +
					" P:" + fmt.Sprintf("%f", k1.High+(k1.High-k1.Low)) +
					" L:" + fmt.Sprintf("%f", k1.High))
				state = 7
			}
		} else if state == 7 {
			if k1.High > k2.High { // means filled
				reset()
			} else {
				//cancel order
				core.Logger.Info("Cancel previous order")
				//place trigger order  in:k1.High loss:k1.Low profit:k1.High+(k1.High-k1.Low)
				core.Logger.Info("[Place Long] Trigger@" + fmt.Sprintf("%f", k1.High) +
					" P:" + fmt.Sprintf("%f", k1.High+(k1.High-k1.Low)) +
					" L:" + fmt.Sprintf("%f", k1.High))
			}
		}

		core.Logger.Info(fmt.Sprintf("====================== state=%d", state))
	}
}

func reset() {
	state = 0
	localHigh = 0
	localLow = 9999999.99
}
