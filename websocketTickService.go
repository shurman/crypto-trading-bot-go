// websocketTickService
package main

import (
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2/futures"
)

func initWsTickService(chanNewKline chan futures.WsKline) {
	log.Println("initWsTickService() Start")

	wsKlineHandler := func(event *futures.WsKlineEvent) {
		log.Println(event.Kline)

		sendIfKlineClosed(event.Kline, chanNewKline)
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

func sendIfKlineClosed(newTick futures.WsKline, chanNewKline chan futures.WsKline) {
	log.Println(newTick.StartTime) //===========
	log.Println(newTick.EndTime)

	chanNewKline <- newTick
}
