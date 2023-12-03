package service

import (
	"crypto-trading-bot-go/core"
)

func BacktestingMode() {
	for _, symbol := range core.Config.Trading.Symbols {
		if core.Config.Trading.Backtesting.Download.Enable {
			DownloadRawHistoryKline(
				symbol,
				core.Config.Trading.Interval,
				core.Config.Trading.Backtesting.Download.StartTime,
				core.Config.Trading.Backtesting.Download.LimitPerDownload)
		}

		LoadRawHistoryKline(symbol, core.Config.Trading.Interval)

		PrintOrderResult(symbol)
		if core.Config.Trading.Backtesting.ExportCsv {
			ExportOrdersResult(symbol)
		}
	}
}
