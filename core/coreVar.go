// coreVar
package core

import (
	"github.com/adshao/go-binance/v2/futures"
)

var (
	NewKline   = make(chan futures.WsKline)
	KLineSlice []futures.WsKline
)
