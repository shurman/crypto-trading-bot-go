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
				v.Fill(currentKline[symbol].CloseTime, false)
				Logger.Debug(fmt.Sprintf("[%s] filled  %+v", v.GetId(), v))
			}
		}
		if v.GetStatus() == core.ORDER_ENTRY {
			if v.GetStopProfitPrice() <= currentKline[symbol].High && v.GetStopProfitPrice() >= currentKline[symbol].Low {
				v.Exit(v.GetStopProfitPrice(), currentKline[symbol].CloseTime, true)
				currentFund[symbol] += v.GetFinalProfit() - v.GetFee()
				Logger.Debug(fmt.Sprintf("[%s] Stop Profit  %+v", v.GetId(), v))

				closedOrdersMap[symbol][v.GetId()] = v
				delete(ordersMap[symbol], v.GetId())
			} else if v.GetStopLossPrice() <= currentKline[symbol].High && v.GetStopLossPrice() >= currentKline[symbol].Low {
				v.Exit(v.GetStopLossPrice(), currentKline[symbol].CloseTime, false)
				currentFund[symbol] += v.GetFinalProfit() - v.GetFee()
				Logger.Debug(fmt.Sprintf("[%s] Stop Loss  %+v", v.GetId(), v))

				closedOrdersMap[symbol][v.GetId()] = v
				delete(ordersMap[symbol], v.GetId())
			}
		}
	}
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
		action = "Cancel & re-Create"
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
		SendSlack(fmt.Sprintf("%s %s %s %.4f@%.4f P:%.4f L:%.4f  (risk:%.3fu) id:%s",
			strategyBO.GetSymbol(),
			action,
			dir.ToString(),
			quantity,
			currentKline[strategyBO.GetSymbol()].Close,
			stopProfit,
			stopLoss,
			(currentKline[strategyBO.GetSymbol()].Close-stopLoss)*quantity*float64(dir),
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
	winLong := 0
	lossLong := 0
	winShort := 0
	lossShort := 0
	profitLong := 0.0
	profitShort := 0.0

	dropDown := 0.0
	maxDropDown := 0.0
	totalFee := 0.0

	for _, v := range closedOrdersMap[symbol] {
		if v.GetDirection() == core.ORDER_LONG {
			countLong++
			if v.GetFinalProfit() > 0 {
				winLong++

				if dropDown < maxDropDown {
					maxDropDown = dropDown
				}
				dropDown = 0
			} else if v.GetFinalProfit() < 0 {
				lossLong++
				dropDown += v.GetFinalProfit()
			}
			profitLong += v.GetFinalProfit()
		} else if v.GetDirection() == core.ORDER_SHORT {
			countShort++
			if v.GetFinalProfit() > 0 {
				winShort++

				if dropDown < maxDropDown {
					maxDropDown = dropDown
				}
				dropDown = 0
			} else if v.GetFinalProfit() < 0 {
				lossShort++
				dropDown += v.GetFinalProfit()
			}
			profitShort += v.GetFinalProfit()
		}

		totalFee += v.GetFee()
	}

	Logger.Warn("=====================================")
	Logger.Warn(fmt.Sprintf("%s Backtesting Result", symbol))
	Logger.Warn("\tLong\t\tShort\t\tTotal")
	Logger.Warn(fmt.Sprintf("Win\t%5d/%5d\t%5d/%5d\t%5d/%5d",
		winLong, winLong+lossLong,
		winShort, winShort+lossShort,
		winLong+winShort, winLong+lossLong+winShort+lossShort))
	winRate := float64(winLong+winShort) / float64(winLong+lossLong+winShort+lossShort)
	Logger.Warn(fmt.Sprintf("Ratio\t=%.3f%%\t=%.3f%%\t=%.3f%% (ev:%.3f)",
		float64(winLong)/float64(winLong+lossLong)*100,
		float64(winShort)/float64(winShort+lossShort)*100,
		winRate*100,
		core.Config.Trading.ProfitLossRatio*winRate-(1-winRate)))
	Logger.Warn(fmt.Sprintf("Profit\t$%-8.2f\t$%-8.2f", profitLong, profitShort))
	Logger.Warn(fmt.Sprintf("Fund\t$%6.2f -> $%6.2f (incl. fee $%5.3f)\t(MDD:%.2f)",
		core.Config.Trading.InitialFund,
		currentFund[symbol],
		totalFee,
		maxDropDown))

	return fmt.Sprintf("%s,%.3f%%,%5d,%.3f,%6.2f,%5.3f\n",
		symbol,
		winRate*100,
		winLong+lossLong+winShort+lossShort,
		core.Config.Trading.ProfitLossRatio*winRate-(1-winRate),
		currentFund[symbol]+totalFee,
		-totalFee)
}

func ExportOrdersResult(symbol string) {
	filename := time.Now().Format("20060102150405") + "_" + symbol + "_orders.csv"
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	var orderKeys []string
	for k := range closedOrdersMap[symbol] {
		orderKeys = append(orderKeys, k)
	}
	sort.Strings(orderKeys)

	for _, k := range orderKeys {
		f.Write([]byte(fmt.Sprintf("%s\n", closedOrdersMap[symbol][k].ToCsv())))
	}

	Logger.Info("Export Order List to " + filename)
}
