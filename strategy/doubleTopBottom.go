// doubleTopBottom
package strategy

import (
	"crypto-trading-bot-go/core"
)

var (
	state     = 0
	borderLow = 0.0
	localHigh = 0.0
)

func doubleTopBottom(notify chan bool) {

	for {
		<-notify

		lenKl := core.GetKlineSliceLen()
		if lenKl < 3 {
			continue
		}

		k1 := core.GetLastKline(1)
		k2 := core.GetLastKline(2)
		k3 := core.GetLastKline(3)

		if state == 0 {
			if k1.High > k2.High && k2.High > k3.High && k1.High/k3.High > 1.01 {

				state = 1
				localHigh = k1.High
				borderLow = k3.Low
			}
		} else if state == 1 {

		}
	}
}
