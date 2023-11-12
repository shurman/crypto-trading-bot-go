// coreVar
package core

import (
	"github.com/adshao/go-binance/v2/futures"
)

var (
	NotifyNewKline = make(chan bool)
	KLineSlice     []futures.WsKline
)
