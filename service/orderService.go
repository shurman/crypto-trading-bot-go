package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
	"log/slog"
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
				slog.Debug(fmt.Sprintf("[%s] filled  %+v", v.GetId(), v))
			}
		} else if v.GetStatus() == core.ORDER_ENTRY {
			if v.GetStopProfitPrice() <= currentKline[symbol].High && v.GetStopProfitPrice() >= currentKline[symbol].Low {
				v.Exit(v.GetStopProfitPrice(), currentKline[symbol].CloseTime)
				CurrentFund[symbol] += v.GetFinalProfit()
				slog.Debug(fmt.Sprintf("[%s] Stop Profit  %+v", v.GetId(), v))
			} else if v.GetStopLossPrice() <= currentKline[symbol].High && v.GetStopLossPrice() >= currentKline[symbol].Low {
				v.Exit(v.GetStopLossPrice(), currentKline[symbol].CloseTime)
				CurrentFund[symbol] += v.GetFinalProfit()
				slog.Debug(fmt.Sprintf("[%s] Stop Loss  %+v", v.GetId(), v))
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
		slog.Warn(strategyBO.GetSymbol() + " Order " + newOrder.GetId() + " not created")
		return
	}

	var action string
	if isReplaced {
		action = "Cancel & re-Create"
	} else {
		action = "Create"
	}

	message := fmt.Sprintf("[%s][%s] %s %s %s %f@%f P:%f L:%f",
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id),
		action,
		strategyBO.GetSymbol(),
		dir.ToString(),
		quantity,
		entry,
		stopProfit,
		stopLoss)
	slog.Debug(message)

	if sendNotify && core.Config.Slack.Enable {
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
		slog.Warn(strategyBO.GetSymbol() + " Order " + newOrder.GetId() + " not created")
		return
	}
	var action string
	if isReplaced {
		action = "Cancel & Entry"
	} else {
		action = "Entry"
	}

	newOrder.Fill(currentKline[strategyBO.GetSymbol()].CloseTime)

	message := fmt.Sprintf("[%s][%s] %s %s %s %f@%f P:%f L:%f",
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id),
		action,
		strategyBO.GetSymbol(),
		dir.ToString(),
		quantity,
		currentKline[strategyBO.GetSymbol()].Close,
		stopProfit,
		stopLoss)
	slog.Debug(message)

	if sendNotify && core.Config.Slack.Enable {
		SendSlack(message)
	}
}

func ExitOrder(
	strategyBO *core.StrategyBO,
	symbol string,
	_id string,
) {
	order, exists := ordersMap[symbol][strategyBO.ToStandardId(_id)]
	if !exists {
		return
	}

	order.Exit(currentKline[symbol].Close, currentKline[symbol].CloseTime)
	CurrentFund[symbol] -= order.GetFinalProfit()
}

func CancelOrder(
	strategyBO *core.StrategyBO,
	symbol string,
	_id string,
	sendNotify bool,
) {
	order, exists := ordersMap[symbol][strategyBO.ToStandardId(_id)]
	if !exists {
		return
	}

	order.Cancel()

	message := fmt.Sprintf("[%s][%s] cancelled",
		currentKline[symbol].CloseTime,
		strategyBO.ToStandardId(_id))
	slog.Debug(message)

	if sendNotify && core.Config.Slack.Enable {
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

	slog.Warn("=====================================")
	slog.Warn(fmt.Sprintf("%s Backtesting Result", symbol))
	slog.Warn("\tLong\t\tShort")
	slog.Warn(fmt.Sprintf("Win\t%5d/%5d\t%5d/%5d", winLong, winLong+lossLong, winShort, winShort+lossShort))
	slog.Warn(fmt.Sprintf("Ratio\t=%3.3f%%\t=%3.3f%%",
		float32(winLong)/float32(winLong+lossLong)*100,
		float32(winShort)/float32(winShort+lossShort)*100))
	slog.Warn(fmt.Sprintf("Profit\t$%5.3f\t$%5.3f", profitLong, profitShort))
	slog.Warn(fmt.Sprintf("Fund\t$%5.3f -> $%5.3f", core.Config.Trading.InitialFund, CurrentFund[symbol]))
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

	slog.Info("Export order list to " + filename)
}
