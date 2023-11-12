// crypto-trading-bot-go project main.go
package main

import (
	"fmt"
	"log"

	//"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

var (
	newTick       = make(chan futures.WsKline)
	strategySlice []Strategy
	kLineSlice    []futures.WsKline
)

func main() {
	go initWsTickService()
	go initStrategyService()

	select {}
}

func initWsTickService() {
	log.Println("initWsTickService() Start")
	wsKlineHandler := func(event *futures.WsKlineEvent) {
		//log.Print("Rev: ")
		log.Println(event.Kline)
		newTick <- event.Kline
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

func initStrategyService() {
	log.Println("initAnalyzeService() Start")
	initCustomStrategies()

	for {
		kLineSlice = append(kLineSlice, getNextClosedKline())

		for _, strategy := range strategySlice {
			strategy.notify <- true
		}
	}
}

type Strategy struct {
	execute func(chan bool)
	notify  chan bool
}

func initCustomStrategies() {
	strategiesSlice := []func(chan bool){doubleTopBottom}

	for _, strategyFunc := range strategiesSlice {
		strategySlice = append(strategySlice, constructStrategy(strategyFunc, make(chan bool)))
	}
}

func constructStrategy(_execute func(chan bool), _notify chan bool) Strategy {
	newStrategy := Strategy{
		execute: _execute,
		notify:  _notify,
	}
	go _execute(_notify)

	return newStrategy
}

func getNextClosedKline() futures.WsKline {
	tick := <-newTick
	log.Println(tick.StartTime) //===========
	log.Println(tick.EndTime)

	nextKline := tick
	return nextKline
}

func doubleTopBottom(notify chan bool) {
	for {
		<-notify
		log.Println("doubleTopBottom")
	}
}
