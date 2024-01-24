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

	currentFund = make(map[string]float64)
)

func init() {
	for _, symbol := range core.Config.Trading.Symbols {
		currentFund[symbol] = core.Config.Trading.InitialFund
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
	currentFund[symbol] += v.GetFinalProfit() - v.GetFee()

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
		SendSlack(fmt.Sprintf("%s %s %s %.4f@%.4f P:%.4f L:%.4f  (risk:%.3fu) id:%s",
			strategyBO.GetSymbol(),
			action,
			dir.ToString(),
			quantity,
			entry,
			stopProfit,
			stopLoss,
			(entry-stopLoss)*quantity*float64(dir),
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
	currentFund[strategyBO.GetSymbol()] -= order.GetFinalProfit()
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
		return currentFund[symbol] * core.Config.Trading.SingleRiskRatio
	} else {
		return core.Config.Trading.InitialFund * core.Config.Trading.SingleRiskRatio
	}
}

func PrintOrderResult(symbol string) string {
	countLong := 0
	countShort := 0
	winLongCount := 0
	lossLongCount := 0
	winShortCount := 0
	lossShortCount := 0
	profitLong := 0.0
	profitShort := 0.0
	lossLong := 0.0
	lossShort := 0.0
	instantWin := 0
	instantLoss := 0

	dropDown := 0.0
	maxDropDown := 0.0
	totalLongFee := 0.0
	totalShortFee := 0.0

	closedOrders := getClosedOrderSlices(symbol)

	for _, v := range closedOrders {
		if v.GetDirection() == core.ORDER_LONG {
			countLong++
			if v.GetFinalProfit() > 0 {
				winLongCount++

				if v.IsInstantFillAndExit() {
					instantWin++
				}

				if dropDown < maxDropDown {
					maxDropDown = dropDown
				}
				dropDown = 0
				profitLong += v.GetFinalProfit() - v.GetFee()

			} else if v.GetFinalProfit() < 0 {
				lossLongCount++
				dropDown += v.GetFinalProfit() - v.GetFee()

				if v.IsInstantFillAndExit() {
					instantLoss++
				}
				lossLong += v.GetFinalProfit() - v.GetFee()
			}
			totalLongFee += v.GetFee()

		} else if v.GetDirection() == core.ORDER_SHORT {
			countShort++
			if v.GetFinalProfit() > 0 {
				winShortCount++

				if v.IsInstantFillAndExit() {
					instantWin++
				}

				if dropDown < maxDropDown {
					maxDropDown = dropDown
				}
				dropDown = 0
				profitShort += v.GetFinalProfit() - v.GetFee()

			} else if v.GetFinalProfit() < 0 {
				lossShortCount++
				dropDown += v.GetFinalProfit() - v.GetFee()

				if v.IsInstantFillAndExit() {
					instantLoss++
				}
				lossShort += v.GetFinalProfit() - v.GetFee()
			}
			totalShortFee += v.GetFee()
		}
	}

	var profitFactor float64
	var recoveryFactor float64
	if lossLong+lossShort == 0 {
		profitFactor = (profitLong + profitShort)
	} else {
		profitFactor = (profitLong + profitShort) / -(lossLong + lossShort)
	}

	if maxDropDown == 0 {
		recoveryFactor = (currentFund[symbol] - core.Config.Trading.InitialFund) / 1
	} else {
		recoveryFactor = (currentFund[symbol] - core.Config.Trading.InitialFund) / -maxDropDown
	}

	Logger.Warn(fmt.Sprintf("%s Backtesting Result", symbol))
	Logger.Warn("\tLong\t\tShort\t\tTotal")
	Logger.Warn(fmt.Sprintf("Win\t%5d/%5d\t%5d/%5d\t%5d/%5d\t(Instant: %2d/%2d) (Still Entry: %d)",
		winLongCount, winLongCount+lossLongCount,
		winShortCount, winShortCount+lossShortCount,
		winLongCount+winShortCount, winLongCount+lossLongCount+winShortCount+lossShortCount,
		instantWin,
		instantWin+instantLoss,
		len(ordersMap[symbol])))
	winRate := float64(winLongCount+winShortCount) / float64(winLongCount+lossLongCount+winShortCount+lossShortCount)
	Logger.Warn(fmt.Sprintf("Ratio\t=%8.3f%%\t=%8.3f%%\t=%8.3f%%\t(ev: %.3f) (P.F.: %.3f) (R.F.: %.3f)",
		float64(winLongCount)/float64(winLongCount+lossLongCount)*100,
		float64(winShortCount)/float64(winShortCount+lossShortCount)*100,
		winRate*100,
		core.Config.Trading.ProfitLossRatio*winRate-(1-winRate),
		profitFactor,
		recoveryFactor))
	Logger.Warn(fmt.Sprintf("Profit\t$%-8.2f\t$%-8.2f\t$%-8.2f", profitLong+lossLong, profitShort+lossShort, profitLong+profitShort+lossLong+lossShort))
	Logger.Warn(fmt.Sprintf("Fee\t$%-8.3f\t$%-8.3f\t$%-8.3f", totalLongFee, totalShortFee, totalLongFee+totalShortFee))
	Logger.Warn(fmt.Sprintf("Fund\t$%6.2f -> $%6.3f (%3.3f%%)\t(MDD:%.2f)",
		core.Config.Trading.InitialFund,
		currentFund[symbol],
		(currentFund[symbol]/core.Config.Trading.InitialFund-1)*100,
		maxDropDown))
	Logger.Warn("=============================================================")

	return fmt.Sprintf("%s,%.3f%%,%5d,%.3f,%6.2f,%5.3f\n",
		symbol,
		winRate*100,
		winLongCount+lossLongCount+winShortCount+lossShortCount,
		profitFactor,
		currentFund[symbol]-core.Config.Trading.InitialFund+(totalLongFee+totalShortFee),
		-(totalLongFee + totalShortFee))
}

func ExportOrdersResult(symbol string) {
	filename := fmt.Sprintf("%s_%s_%s_%.2f_orders.csv",
		time.Now().Format("20060102150405"),
		symbol,
		core.Config.Trading.Interval,
		core.Config.Trading.ProfitLossRatio)
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	closedOrders := getClosedOrderSlices(symbol)
	for _, v := range closedOrders {
		f.Write([]byte(fmt.Sprintf("%s\n", v.ToCsv())))
	}

	Logger.Info("Export Order List to " + filename)
}

func ExportPerformanceChart(symbol string) {
	filename := fmt.Sprintf("%s_%s_%.2f.png",
		symbol,
		core.Config.Trading.Interval,
		core.Config.Trading.ProfitLossRatio)
	closedOrders := getClosedOrderSlices(symbol)

	var xValues []time.Time
	var yValues []float64
	currentFund := core.Config.Trading.InitialFund

	for _, v := range closedOrders {
		currentFund += v.GetFinalProfit() - v.GetFee()
		xValues = append(xValues, v.GetCreateTime())
		yValues = append(yValues, currentFund)
	}

	ExportChart(xValues, yValues, symbol+" "+core.Config.Trading.Interval, filename)
	Logger.Info("Export Performance Chart to " + filename)
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
