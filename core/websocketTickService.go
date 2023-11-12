// websocketTickService
package core

import (
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2/futures"
)

func InitWsTickService() {
	log.Println("initWsTickService() Start")

	wsKlineHandler := func(event *futures.WsKlineEvent) {
		log.Println(event.Kline)

		sendIfKlineClosed(event.Kline)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	log.Println(wsKlineHandler)

	doneC, _, err := futures.WsKlineServe("ETHUSDT", "1m", wsKlineHandler, errHandler) //WsCombinedKlineServe for multiple
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}

func sendIfKlineClosed(newTick futures.WsKline) {
	log.Println(newTick.StartTime) //===========
	log.Println(newTick.EndTime)

	KLineSlice = append(KLineSlice, newTick)

	NotifyNewKline <- true
}
