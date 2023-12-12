package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
	"os"
	"sort"
	"time"
)

var (
	ordersMap    = make(map[string](map[string]*core.OrderBO))
	currentKline = make(map[string]*core.Kline)

	CurrentFund = make(map[string]float64)
)

func init() {
	for _, symbol := range core.Config.Trading.Symbols {
		CurrentFund[symbol] = core.Config.Trading.InitialFund
	}
}

func CheckOrderFilled(symbol string, k *core.Kline) {
	currentKline[symbol] = k

	for _, v := range ordersMap[symbol] {
		if v.GetStatus() == core.ORDER_OPEN {
			if v.GetEntryPrice() <= currentKline[symbol].High && v.GetEntryPrice() >= currentKline[symbol].Low {
				v.Fill(currentKline[symbol].CloseTime)
				Logger.Debug(fmt.Sprintf("[%s] filled  %+v", v.GetId(), v))
			}
		} else if v.GetStatus() == core.ORDER_ENTRY {
			if v.GetStopProfitPrice() <= currentKline[symbol].High && v.GetStopProfitPrice() >= currentKline[symbol].Low {
				v.Exit(v.GetStopProfitPrice(), currentKline[symbol].CloseTime)
				CurrentFund[symbol] += v.GetFinalProfit()
				Logger.Debug(fmt.Sprintf("[%s] Stop Profit  %+v", v.GetId(), v))
			} else if v.GetStopLossPrice() <= currentKline[symbol].High && v.GetStopLossPrice() >= currentKline[symbol].Low {
				v.Exit(v.GetStopLossPrice(), currentKline[symbol].CloseTime)
				CurrentFund[symbol] += v.GetFinalProfit()
				Logger.Debug(fmt.Sprintf("[%s] Stop Loss  %+v", v.GetId(), v))
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

	message := fmt.Sprintf("[%s] %s | %s %s %s %.2f@%.2f P:%.2f L:%.2f  (risk:%.2fu)",
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
		SendSlack(message)
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

	newOrder.Fill(currentKline[strategyBO.GetSymbol()].CloseTime)

	message := fmt.Sprintf("[%s] %s | %s %s %s %.2f@%.2f P:%.2f L:%.2f  (risk:%.2fu)",
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
		SendSlack(message)
	}
}

func ExitOrder(
	strategyBO *core.StrategyBO,
	_id string,
) {
	order, exists := ordersMap[strategyBO.GetSymbol()][strategyBO.ToStandardId(_id)]
	if !exists {
		return
	}

	order.Exit(currentKline[strategyBO.GetSymbol()].Close, currentKline[strategyBO.GetSymbol()].CloseTime)
	CurrentFund[strategyBO.GetSymbol()] -= order.GetFinalProfit()
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

func GetOrderStatus(strategyBO *core.StrategyBO, id string) core.OrderStatus {
	order, exists := ordersMap[strategyBO.GetSymbol()][strategyBO.ToStandardId(id)]

	if exists {
		return order.GetStatus()
	}
	return core.ORDER_UNKWN
}

func orderPut(symbol string, id string, newOrder *core.OrderBO) (success bool, isReplaced bool) {
	if _, exists := ordersMap[symbol]; !exists {
		ordersMap[symbol] = make(map[string]*core.OrderBO)
	} else if order, exists := ordersMap[symbol][id]; exists {
		if order.GetStatus() == core.ORDER_OPEN {
			ordersMap[symbol][id] = newOrder
			return true, true
		} else {
			return false, false
		}
	}

	ordersMap[symbol][id] = newOrder
	return true, false
}

func PrintOrderResult(symbol string) {
	countLong := 0
	countShort := 0
	winLong := 0
	lossLong := 0
	winShort := 0
	lossShort := 0
	profitLong := 0.0
	profitShort := 0.0

	for _, v := range ordersMap[symbol] {
		if v.GetDirection() == core.ORDER_LONG {
			countLong++
			if v.GetFinalProfit() > 0 {
				winLong++
			} else {
				lossLong++
			}
			profitLong += v.GetFinalProfit()
		} else if v.GetDirection() == core.ORDER_SHORT {
			countShort++
			if v.GetFinalProfit() > 0 {
				winShort++
			} else {
				lossShort++
			}
			profitShort += v.GetFinalProfit()
		}
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
	Logger.Warn(fmt.Sprintf("Fund\t$%6.2f -> $%6.2f", core.Config.Trading.InitialFund, CurrentFund[symbol]))
}

func ExportOrdersResult(symbol string) {
	filename := time.Now().Format("20060102150405") + "_" + symbol + "_report.csv"
	f, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	var orderKeys []string
	for k := range ordersMap[symbol] {
		orderKeys = append(orderKeys, k)
	}
	sort.Strings(orderKeys)

	for _, k := range orderKeys {
		f.Write([]byte(fmt.Sprintf("%s\n", ordersMap[symbol][k].ToCsv())))
	}

	Logger.Info("Export order list to " + filename)
}
