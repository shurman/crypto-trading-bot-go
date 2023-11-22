// kline
package core

import (
	"fmt"
)

type Kline struct {
	StartTime int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	IsNew     bool
}

func (k *Kline) ToString() string {
	return fmt.Sprintf("{s_time: %d, o:%f, h:%f, l:%f, c:%f}", k.StartTime, k.Open, k.High, k.Low, k.Close)
}
