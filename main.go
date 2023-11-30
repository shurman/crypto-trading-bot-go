// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//DTB reset after long waiting / state 1,-1 rule rework / phase 2
//quantity based on Compound Interest
//multiple coins

func main() {
	service.InitStrategyService()

	if service.Config.Trading.Mode == "indicator" {
		indicatorMode()
	} else {
		backtestingMode()
	}
}

func indicatorMode() {
	service.LoadHistoryKline()
	go service.InitWsTickService()
	select {}
}

func backtestingMode() {
	if service.Config.Trading.Backtesting.Download.Enable {
		service.DownloadRawHistoryKline(
			service.Config.Trading.Symbol,
			service.Config.Trading.Interval,
			service.Config.Trading.Backtesting.Download.StartTime,
			service.Config.Trading.Backtesting.Download.LimitPerDownload)
	}
	service.LoadRawHistoryKline(service.Config.Trading.Symbol, service.Config.Trading.Interval)
	service.OutputOrdersResult()
}
