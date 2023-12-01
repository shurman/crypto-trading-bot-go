// crypto-trading-bot-go project main.go
package main

import (
	"crypto-trading-bot-go/core"
	"crypto-trading-bot-go/service"
	_ "crypto-trading-bot-go/strategy"
)

//TODO
//quantity based on Compound Interest
//print order result after backtesting
//try 1h
//multiple coins
//DTB reset after long waiting / state 1,-1 rule rework / phase 2

func main() {
	service.InitStrategyService()

	if core.Config.Trading.Mode == "indicator" {
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
	if core.Config.Trading.Backtesting.Download.Enable {
		service.DownloadRawHistoryKline(
			core.Config.Trading.Symbol,
			core.Config.Trading.Interval,
			core.Config.Trading.Backtesting.Download.StartTime,
			core.Config.Trading.Backtesting.Download.LimitPerDownload)
	}
	service.LoadRawHistoryKline(core.Config.Trading.Symbol, core.Config.Trading.Interval)

	if core.Config.Trading.Backtesting.ExportCsv {
		service.OutputOrdersResult()
	}
}
