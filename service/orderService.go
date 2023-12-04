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

	CurrentFund = core.Config.Trading.InitialFund
)

func CheckOrderFilled(symbol string, k *core.Kline) {
	currentKline[symbol] = k

	for _, v := range ordersMap[symbol] {
		if v.GetStatus() == core.ORDER_OPEN {
			if v.GetEntryPrice() <= currentKline[symbol].High && v.GetEntryPrice() >= currentKline[symbol].Low {
				v.Fill(currentKline[symbol].CloseTime)
				slog.Info(fmt.Sprintf("[%s] filled  %+v", v.GetId(), v))
			}
		} else if v.GetStatus() == core.ORDER_ENTRY {
			if v.GetStopProfitPrice() <= currentKline[symbol].High && v.GetStopProfitPrice() >= currentKline[symbol].Low {
				v.Exit(v.GetStopProfitPrice(), currentKline[symbol].CloseTime)
				CurrentFund += v.GetFinalProfit()
				slog.Info(fmt.Sprintf("[%s] Stop Profit  %+v", v.GetId(), v))
			} else if v.GetStopLossPrice() <= currentKline[symbol].High && v.GetStopLossPrice() >= currentKline[symbol].Low {
				v.Exit(v.GetStopLossPrice(), currentKline[symbol].CloseTime)
				CurrentFund += v.GetFinalProfit()
				slog.Info(fmt.Sprintf("[%s] Stop Loss  %+v", v.GetId(), v))
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

	if !orderPut(strategyBO.GetSymbol(), newOrder.GetId(), newOrder) {
		slog.Warn(strategyBO.GetSymbol() + " Order " + newOrder.GetId() + " not created")
		return
	}

	slog.Info(fmt.Sprintf("[%s][%s] create %s %s %f@%f P:%f L:%f",
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id),
		strategyBO.GetSymbol(),
		dir.ToString(),
		quantity,
		entry,
		stopProfit,
		stopLoss))

	if sendNotify {
		//send slack
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

	if !orderPut(strategyBO.GetSymbol(), newOrder.GetId(), newOrder) {
		slog.Warn("Order " + newOrder.GetId() + " not created")
		return
	}

	newOrder.Fill(currentKline[strategyBO.GetSymbol()].CloseTime)

	slog.Info(fmt.Sprintf("[%s][%s] entry %s %f@%f P:%f L:%f",
		currentKline[strategyBO.GetSymbol()].CloseTime,
		strategyBO.ToStandardId(_id),
		dir.ToString(),
		quantity,
		currentKline[strategyBO.GetSymbol()].Close,
		stopProfit,
		stopLoss))

	if sendNotify {
		//send slack
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
	CurrentFund -= order.GetFinalProfit()
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

	slog.Info(fmt.Sprintf("[%s][%s] cancelled",
		currentKline[symbol].CloseTime,
		strategyBO.ToStandardId(_id),
	))

	if sendNotify {
		//send slack
	}
}

func GetOrderStatus(strategyBO *core.StrategyBO, id string) core.OrderStatus {
	order, exists := ordersMap[strategyBO.GetSymbol()][strategyBO.ToStandardId(id)]

	if exists {
		return order.GetStatus()
	}
	return core.ORDER_UNKWN
}

func orderPut(symbol string, id string, newOrder *core.OrderBO) bool {
	if order, exists := ordersMap[symbol][id]; exists {
		if order.GetStatus() == core.ORDER_OPEN {
			ordersMap[symbol][id] = newOrder
			return true
		} else {
			return false
		}
	}

	ordersMap[symbol][id] = newOrder
	return true
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

	slog.Warn(fmt.Sprintf("%s Backtesting Result", symbol))
	slog.Warn("\tLong\t\tShort")
	slog.Warn(fmt.Sprintf("Win\t%5d/%5d\t%5d/%5d", winLong, winLong+lossLong, winShort, winShort+lossShort))
	slog.Warn(fmt.Sprintf("Ratio\t=%3.3f%%\t=%3.3f%%",
		float32(winLong)/float32(winLong+lossLong)*100,
		float32(winShort)/float32(winShort+lossShort)*100))
	slog.Warn(fmt.Sprintf("Profit\t$%5.3f\t$%5.3f", profitLong, profitShort))
	slog.Warn(fmt.Sprintf("Fund\t$%5.3f -> $%5.3f", core.Config.Trading.InitialFund, CurrentFund))
}

func ExportOrdersResult(symbol string) {
	f, _ := os.OpenFile(time.Now().Format("20060102150405")+"_"+symbol+"_report.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	var orderKeys []string
	for k := range ordersMap[symbol] {
		orderKeys = append(orderKeys, k)
	}
	sort.Strings(orderKeys)

	for _, k := range orderKeys {
		f.Write([]byte(fmt.Sprintf("%s\n", ordersMap[symbol][k].ToCsv())))
	}
}
