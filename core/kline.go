// kline
package core

import (
	"fmt"
	"time"
)

type Kline struct {
	StartTime time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	CloseTime time.Time
	IsNew     bool
}

func (k *Kline) ToString() string {
	return fmt.Sprintf("{s_time: %s, o:%f, h:%f, l:%f, c:%f}", k.StartTime, k.Open, k.High, k.Low, k.Close)
}
