// kline
package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
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
	return fmt.Sprintf("{s_time: %s, o:%.2f, h:%.2f, l:%.2f, c:%.2f}", k.StartTime, k.Open, k.High, k.Low, k.Close)
}

func ConvertToKlineFromWsKline(tick *futures.WsKline) Kline {
	return Kline{
		StartTime: time.Unix(tick.StartTime/1000, 0),
		Open:      parseFloat(tick.Open),
		High:      parseFloat(tick.High),
		Low:       parseFloat(tick.Low),
		Close:     parseFloat(tick.Close),
		CloseTime: time.Unix(tick.EndTime/1000, 0),
		IsNew:     true,
	}
}

func ConvertKlineFromFuturesKline(tick *futures.Kline) Kline {
	return Kline{
		StartTime: time.Unix(tick.OpenTime/1000, 0),
		Open:      parseFloat(tick.Open),
		High:      parseFloat(tick.High),
		Low:       parseFloat(tick.Low),
		Close:     parseFloat(tick.Close),
		CloseTime: time.Unix(tick.CloseTime/1000, 0),
		IsNew:     false,
	}
}

func parseFloat(str string) float64 {
	result, _ := strconv.ParseFloat(str, 64)

	return result
}
