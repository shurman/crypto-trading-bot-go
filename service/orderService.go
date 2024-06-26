package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
	"os"
	"sort"
	"time"
)

var (
	ordersMap       = make(map[string](map[string]*core.OrderBO))
	closedOrdersMap = make(map[string](map[string]*core.OrderBO))
	currentKline    = make(map[string]*core.Kline)

	currentFundMap = make(map[string]float64)
)

func init() {
	for _, symbol := range core.Config.Trading.Symbols {
		currentFundMap[symbol] = core.Config.Trading.InitialFund
	}
}

func CheckOrderFilled(symbol string, k *core.Kline) {
	currentKline[symbol] = k

	for _, v := range ordersMap[symbol] {
		if v.GetStatus() == core.ORDER_OPEN {
			if v.GetEntryPrice() <= currentKline[symbol].High && v.GetEntryPrice() >= currentKline[symbol].Low {
				v.Fill(currentKline[symbol].CloseTime, true)
				Logger.Debug(fmt.Sprintf("[%s] filled  %+v", v.GetId(), v))
			}
		}
		if v.GetStatus() == core.ORDER_ENTRY {
			if v.GetStopProfitPrice() <= currentKline[symbol].High && v.GetStopProfitPrice() >= currentKline[symbol].Low {
				ExitingOrder(v, v.GetStopProfitPrice(), symbol, true, false)

			} else if v.GetStopLossPrice() <= currentKline[symbol].High && v.GetStopLossPrice() >= currentKline[symbol].Low {
				ExitingOrder(v, v.GetStopLossPrice(), symbol, false, true)

			} else if v.GetDirection() == core.ORDER_LONG && v.GetStopProfitPrice() < currentKline[symbol].Low {
				ExitingOrder(v, currentKline[symbol].Close, symbol, true, true)

			} else if v.GetDirection() == core.ORDER_SHORT && v.GetStopProfitPrice() > currentKline[symbol].High {
				ExitingOrder(v, currentKline[symbol].Close, symbol, true, true)

			} else if v.GetDirection() == core.ORDER_LONG && v.GetStopLossPrice() > currentKline[symbol].High {
				ExitingOrder(v, currentKline[symbol].Close, symbol, false, true)

			} else if v.GetDirection() == core.ORDER_SHORT && v.GetStopLossPrice() < currentKline[symbol].Low {
				ExitingOrder(v, currentKline[symbol].Close, symbol, false, true)
			}
		}
	}
}

func ExitingOrder(v *core.OrderBO, exitPrice float64, symbol string, isTakeProfit bool, isTaker bool) {
	v.Exit(exitPrice, currentKline[symbol].CloseTime, isTaker)
	currentFundMap[symbol] += v.GetFinalProfit() - v.GetFee()

	if isTakeProfit {
		Logger.Debug(fmt.Sprintf("[%s] Take Profit  %+v", v.GetId(), v))
	} else {
		Logger.Debug(fmt.Sprintf("[%s] Stop Loss  %+v", v.GetId(), v))
	}

	closedOrdersMap[symbol][v.GetId()] = v
	delete(ordersMap[symbol], v.GetId())
}

func CreateOrder(
	strategyBO *core.StrategyBO,
	_id string,
	dir core.OrderDirection,
	quantity float64,
	entry float64,
	stopProfit float64,
	stopLoss float64,
	sendNotify bool,
) {

	newOrder := core.ConstructOrderBO(
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id),
		core.ORDER_OPEN,
		dir,
		quantity,
		entry,
		stopProfit,
		stopLoss,
	)

	success, isReplaced := orderPut(strategyBO.GetSymbol(), newOrder.GetId(), newOrder)
	if !success {
		Logger.Warn(strategyBO.GetSymbol() + " Order " + newOrder.GetId() + " not created")
		return
	}

	var action string
	if isReplaced {
		action = "re-Create"
	} else {
		action = "Create"
	}

	message := fmt.Sprintf("[%s] %s | %s %s %s %.4f@%.4f P:%.4f L:%.4f  (risk:%.3fu)",
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id),
		action,
		strategyBO.GetSymbol(),
		dir.ToString(),
		quantity,
		entry,
		stopProfit,
		stopLoss,
		(entry-stopLoss)*quantity*float64(dir))
	if sendNotify {
		Logger.Info(message)
	} else {
		Logger.Debug(message)
	}

	if sendNotify {
		SendSlack(fmt.Sprintf("%s %s %s@%.4f P:%.4f L:%.4f  %s",
			strategyBO.GetSymbol(),
			action,
			dir.ToString(),
			entry,
			stopProfit,
			stopLoss,
			strategyBO.ToStandardId(_id),
		))
	}
}

func CreateMarketOrder(strategyBO *core.StrategyBO,
	_id string,
	dir core.OrderDirection,
	quantity float64,
	stopProfit float64,
	stopLoss float64,
	sendNotify bool,
) {
	newOrder := core.ConstructOrderBO(
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id),
		core.ORDER_OPEN,
		dir,
		quantity,
		currentKline[strategyBO.GetSymbol()].Close,
		stopProfit,
		stopLoss,
	)

	success, isReplaced := orderPut(strategyBO.GetSymbol(), newOrder.GetId(), newOrder)
	if !success {
		Logger.Warn(strategyBO.GetSymbol() + " Order " + newOrder.GetId() + " not created")
		return
	}
	var action string
	if isReplaced {
		action = "Cancel & Entry"
	} else {
		action = "Entry"
	}

	newOrder.Fill(currentKline[strategyBO.GetSymbol()].CloseTime, true)

	message := fmt.Sprintf("[%s] %s | %s %s %s %.4f@%.4f P:%.4f L:%.4f  (risk:%.3fu)",
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id),
		action,
		strategyBO.GetSymbol(),
		dir.ToString(),
		quantity,
		currentKline[strategyBO.GetSymbol()].Close,
		stopProfit,
		stopLoss,
		(currentKline[strategyBO.GetSymbol()].Close-stopLoss)*quantity*float64(dir))
	if sendNotify {
		Logger.Info(message)
	} else {
		Logger.Debug(message)
	}

	if sendNotify {
		SendSlack(fmt.Sprintf("%s %s %s@%.4f P:%.4f L:%.4f  %s",
			strategyBO.GetSymbol(),
			action,
			dir.ToString(),
			currentKline[strategyBO.GetSymbol()].Close,
			stopProfit,
			stopLoss,
			strategyBO.ToStandardId(_id),
		))
	}
}

func ExitOrder(
	strategyBO *core.StrategyBO,
	_id string,
	isTaker bool,
) {
	order, exists := ordersMap[strategyBO.GetSymbol()][strategyBO.ToStandardId(_id)]
	if !exists {
		return
	}

	order.Exit(currentKline[strategyBO.GetSymbol()].Close, currentKline[strategyBO.GetSymbol()].CloseTime, isTaker)
	currentFundMap[strategyBO.GetSymbol()] -= order.GetFinalProfit()
}

func CancelOrder(
	strategyBO *core.StrategyBO,
	_id string,
	sendNotify bool,
) {
	order, exists := ordersMap[strategyBO.GetSymbol()][strategyBO.ToStandardId(_id)]
	if !exists {
		return
	}

	order.Cancel()

	message := fmt.Sprintf("[%s] %s | cancelled",
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id))
	if sendNotify {
		Logger.Info(message)
	} else {
		Logger.Debug(message)
	}

	if sendNotify {
		SendSlack(message)
	}
}

func GetOrder(strategyBO *core.StrategyBO, id string) *core.OrderBO {
	order, exists := ordersMap[strategyBO.GetSymbol()][strategyBO.ToStandardId(id)]

	if exists {
		return order
	}

	order, exists = closedOrdersMap[strategyBO.GetSymbol()][strategyBO.ToStandardId(id)]
	if exists {
		return order
	}
	return nil
}

func orderPut(symbol string, id string, newOrder *core.OrderBO) (success bool, isReplaced bool) {
	if _, exists := ordersMap[symbol]; !exists {
		ordersMap[symbol] = make(map[string]*core.OrderBO)
		closedOrdersMap[symbol] = make(map[string]*core.OrderBO)
	} else if order, exists := ordersMap[symbol][id]; exists {
		if order.GetStatus() == core.ORDER_OPEN || order.GetStatus() == core.ORDER_CANCEL {
			ordersMap[symbol][id] = newOrder
			return true, true
		} else {
			return false, false
		}
	}

	ordersMap[symbol][id] = newOrder
	return true, false
}

func GetRiskPerTrade(symbol string) float64 {
	if core.Config.Trading.EnableAccumulated {
		return currentFundMap[symbol] * core.Config.Trading.SingleRiskRatio
	} else {
		return core.Config.Trading.InitialFund * core.Config.Trading.SingleRiskRatio
	}
}

func getClosedOrderSlices(symbol string) (closedOrders []*core.OrderBO) {
	for _, v := range closedOrdersMap[symbol] {
		closedOrders = append(closedOrders, v)
	}
	sort.Slice(closedOrders, func(i int, j int) bool {
		return closedOrders[i].GetCreateTime().Before(closedOrders[j].GetCreateTime())
	})

	return
}

func PrintSymbolResult(symbol string) string {
	bo := core.CalculateOrdersResult(getClosedOrderSlices(symbol))
	profitFactor, recoveryFactor := core.GetPfRf(bo.ProfitLong+bo.ProfitShort, bo.LossLong+bo.LossShort, currentFundMap[symbol]-core.Config.Trading.InitialFund, bo.MaxDropDown)

	Logger.Warn(fmt.Sprintf("%s Backtesting Result", symbol))
	Logger.Warn("\tLong\t\tShort\t\tTotal")
	Logger.Warn(fmt.Sprintf("Win\t%5d/%5d\t%5d/%5d\t%5d/%5d\t(Instant: %2d/%2d) (Still Entry: %d)",
		bo.WinLongCount, bo.WinLongCount+bo.LossLongCount,
		bo.WinShortCount, bo.WinShortCount+bo.LossShortCount,
		bo.WinLongCount+bo.WinShortCount, bo.WinLongCount+bo.LossLongCount+bo.WinShortCount+bo.LossShortCount,
		bo.InstantWin,
		bo.InstantWin+bo.InstantLoss,
		len(ordersMap[symbol])))
	winRate := float64(bo.WinLongCount+bo.WinShortCount) / float64(bo.WinLongCount+bo.LossLongCount+bo.WinShortCount+bo.LossShortCount)
	Logger.Warn(fmt.Sprintf("Ratio\t=%8.3f%%\t=%8.3f%%\t=%8.3f%%\t(ev: %.3f) (P.F.: %.3f) (R.F.: %.3f)",
		float64(bo.WinLongCount)/float64(bo.WinLongCount+bo.LossLongCount)*100,
		float64(bo.WinShortCount)/float64(bo.WinShortCount+bo.LossShortCount)*100,
		winRate*100,
		core.Config.Trading.ProfitLossRatio*winRate-(1-winRate),
		profitFactor,
		recoveryFactor))
	Logger.Warn(fmt.Sprintf("Profit\t$%-8.3f\t$%-8.3f\t$%-8.3f", bo.ProfitLong+bo.LossLong, bo.ProfitShort+bo.LossShort, bo.ProfitLong+bo.ProfitShort+bo.LossLong+bo.LossShort))
	Logger.Warn(fmt.Sprintf("Fee\t$%-8.3f\t$%-8.3f\t$%-8.3f", bo.TotalLongFee, bo.TotalShortFee, bo.TotalLongFee+bo.TotalShortFee))
	Logger.Warn(fmt.Sprintf("Fund\t$%6.2f -> $%6.3f (%3.3f%%)\t(MDD:%.2f)",
		core.Config.Trading.InitialFund,
		currentFundMap[symbol],
		(currentFundMap[symbol]/core.Config.Trading.InitialFund-1)*100,
		bo.MaxDropDown))
	Logger.Warn("=============================================================")

	return fmt.Sprintf("%s,%.3f%%,%5d,%.3f,%6.2f,%5.3f\n",
		symbol,
		winRate*100,
		bo.WinLongCount+bo.LossLongCount+bo.WinShortCount+bo.LossShortCount,
		profitFactor,
		currentFundMap[symbol]-core.Config.Trading.InitialFund+(bo.TotalLongFee+bo.TotalShortFee),
		-(bo.TotalLongFee + bo.TotalShortFee))
}

func ExportSymbolResult(symbol string) {
	if !core.Config.Trading.Backtesting.Export.Orders && core.Config.Trading.Backtesting.Export.Chart {
		return
	}

	closedOrders := getClosedOrderSlices(symbol)
	var rawOutputOrders = ""
	var xValues []time.Time
	var yValues []float64
	currentFund := core.Config.Trading.InitialFund

	for _, v := range closedOrders {
		if core.Config.Trading.Backtesting.Export.Orders {
			rawOutputOrders += v.ToCsv() + "\n"
		}
		if core.Config.Trading.Backtesting.Export.Chart {
			currentFund += v.GetFinalProfit() - v.GetFee()
			xValues = append(xValues, v.GetCreateTime())
			yValues = append(yValues, currentFund)
		}
	}

	if core.Config.Trading.Backtesting.Export.Orders {
		orderFilename := fmt.Sprintf("%s%s_%s_%s_%.2f_orders.csv",
			core.ReportFilePath,
			time.Now().Format("20060102150405"),
			symbol,
			core.Config.Trading.Interval,
			core.Config.Trading.ProfitLossRatio)
		f, _ := os.OpenFile(orderFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.Write([]byte(rawOutputOrders))
		Logger.Info("Export Order List to " + orderFilename)
	}
	if core.Config.Trading.Backtesting.Export.Chart {
		chartFilename := fmt.Sprintf("%s%s_%s_%.2f",
			core.ChartFilePath,
			symbol,
			core.Config.Trading.Interval,
			core.Config.Trading.ProfitLossRatio)
		ExportChart(xValues, yValues, symbol+" "+core.Config.Trading.Interval, chartFilename)
		Logger.Info("Export Performance Chart to " + chartFilename)
	}
}

func ExportOverallPerformanceChart() {
	var closedOrders []*core.OrderBO
	for _, symbol := range core.Config.Trading.Symbols {
		closedOrders = append(closedOrders, getClosedOrderSlices(symbol)...)
	}

	sort.Slice(closedOrders, func(i int, j int) bool {
		return closedOrders[i].GetCreateTime().Before(closedOrders[j].GetCreateTime())
	})

	var xValues []time.Time
	var yValues []float64
	initalFund := core.Config.Trading.InitialFund * float64(len(core.Config.Trading.Symbols))
	currentFund := initalFund

	bo := core.CalculateOrdersResult(closedOrders)

	for _, v := range closedOrders {
		currentFund += v.GetFinalProfit() - v.GetFee()

		if core.Config.Trading.Backtesting.OverallChart {
			xValues = append(xValues, v.GetCreateTime())
			yValues = append(yValues, currentFund)
		}
	}

	profitFactor, recoveryFactor := core.GetPfRf(bo.ProfitLong+bo.ProfitShort, bo.LossLong+bo.LossShort, currentFund-core.Config.Trading.InitialFund, bo.MaxDropDown)

	Logger.Warn("Overall Backtesting Result")
	Logger.Warn("\tLong\t\tShort\t\tTotal")
	Logger.Warn(fmt.Sprintf("Win\t%5d/%5d\t%5d/%5d\t%5d/%5d\t(Instant: %2d/%2d)",
		bo.WinLongCount, bo.WinLongCount+bo.LossLongCount,
		bo.WinShortCount, bo.WinShortCount+bo.LossShortCount,
		bo.WinLongCount+bo.WinShortCount, bo.WinLongCount+bo.LossLongCount+bo.WinShortCount+bo.LossShortCount,
		bo.InstantWin,
		bo.InstantWin+bo.InstantLoss))
	winRate := float64(bo.WinLongCount+bo.WinShortCount) / float64(bo.WinLongCount+bo.LossLongCount+bo.WinShortCount+bo.LossShortCount)
	Logger.Warn(fmt.Sprintf("Ratio\t=%8.3f%%\t=%8.3f%%\t=%8.3f%%\t(ev: %.3f) (P.F.: %.3f) (R.F.: %.3f)",
		float64(bo.WinLongCount)/float64(bo.WinLongCount+bo.LossLongCount)*100,
		float64(bo.WinShortCount)/float64(bo.WinShortCount+bo.LossShortCount)*100,
		winRate*100,
		core.Config.Trading.ProfitLossRatio*winRate-(1-winRate),
		profitFactor,
		recoveryFactor))
	Logger.Warn(fmt.Sprintf("Profit\t$%-8.3f\t$%-8.3f\t$%-8.3f", bo.ProfitLong+bo.LossLong, bo.ProfitShort+bo.LossShort, bo.ProfitLong+bo.ProfitShort+bo.LossLong+bo.LossShort))
	Logger.Warn(fmt.Sprintf("Fee\t$%-8.3f\t$%-8.3f\t$%-8.3f", bo.TotalLongFee, bo.TotalShortFee, bo.TotalLongFee+bo.TotalShortFee))
	Logger.Warn(fmt.Sprintf("Fund\t$%6.2f -> $%6.3f (%3.3f%%)\t(MDD:%.2f)",
		initalFund,
		currentFund,
		(currentFund/initalFund-1)*100,
		bo.MaxDropDown))
	Logger.Warn("=============================================================")

	if core.Config.Trading.Backtesting.OverallChart {
		chartFilename := fmt.Sprintf("%soverall_%s_%.2f",
			core.ChartFilePath,
			core.Config.Trading.Interval,
			core.Config.Trading.ProfitLossRatio)
		chartTitle := fmt.Sprintf("Overall %s %.2f",
			core.Config.Trading.Interval,
			core.Config.Trading.ProfitLossRatio)
		ExportChart(xValues, yValues, chartTitle, chartFilename)
		Logger.Info("Export Performance Chart to " + chartFilename)
	}
}
