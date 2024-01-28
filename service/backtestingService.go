package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
	"os"
	"time"
)

func BacktestingMode() {
	filename := fmt.Sprintf("%s%s_%s_%.2f_reports.csv",
		core.ReportFilePath,
		time.Now().Format("20060102150405"),
		core.Config.Trading.Interval,
		core.Config.Trading.ProfitLossRatio)

	for _, symbol := range core.Config.Trading.Symbols {
		if core.Config.Trading.Backtesting.Download.Enable {
			DownloadRawHistoryKline(
				symbol,
				core.Config.Trading.Interval,
				core.Config.Trading.Backtesting.Download.StartTime,
				core.Config.Trading.Backtesting.Download.LimitPerDownload)
		}

		LoadRawHistoryKline(symbol, core.Config.Trading.Interval)

		summary := PrintOrderResult(symbol)
		if core.Config.Trading.Backtesting.Export.Reports {
			f, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			f.Write([]byte(summary))
		}

		ExportSymbolResult(symbol)
	}

	if core.Config.Trading.Backtesting.Export.Reports {
		Logger.Info("Export Summary Report to " + filename)
	}

	ExportOverallPerformanceChart()
}
