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
	orderMap     = make(map[string]*core.OrderBO)
	currentKline *core.Kline
)

func SetCurrentKline(k *core.Kline) {
	currentKline = k
}

func CheckOrderFilled() {
	for _, v := range orderMap {
		if v.GetStatus() == core.ORDER_OPEN {
			if v.GetEntryPrice() <= currentKline.High && v.GetEntryPrice() >= currentKline.Low {
				v.Fill(currentKline.CloseTime)
				slog.Info(fmt.Sprintf("[%s] filled  %+v", v.GetId(), v))
			}
		} else if v.GetStatus() == core.ORDER_ENTRY {
			if v.GetStopProfitPrice() <= currentKline.High && v.GetStopProfitPrice() >= currentKline.Low {
				v.Exit(v.GetStopProfitPrice(), currentKline.CloseTime)
				slog.Info(fmt.Sprintf("[%s] Stop Profit  %+v", v.GetId(), v))
			} else if v.GetStopLossPrice() <= currentKline.High && v.GetStopLossPrice() >= currentKline.Low {
				v.Exit(v.GetStopLossPrice(), currentKline.CloseTime)
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
		currentKline.CloseTime,
		strategyBO.ToStandardId(_id),
		core.ORDER_OPEN,
		dir,
		quantity,
		entry,
		stopProfit,
		stopLoss,
	)

	if !orderPut(newOrder.GetId(), newOrder) {
		slog.Warn("Order " + newOrder.GetId() + " not created")
		return
	}

	slog.Info(fmt.Sprintf("[%s][%s] create %s %f@%f P:%f L:%f",
		currentKline.CloseTime,
		strategyBO.ToStandardId(_id),
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
		currentKline.CloseTime,
		strategyBO.ToStandardId(_id),
		core.ORDER_OPEN,
		dir,
		quantity,
		currentKline.Close,
		stopProfit,
		stopLoss,
	)

	if !orderPut(newOrder.GetId(), newOrder) {
		slog.Warn("Order " + newOrder.GetId() + " not created")
		return
	}

	newOrder.Fill(currentKline.CloseTime)

	slog.Info(fmt.Sprintf("[%s][%s] entry %s %f@%f P:%f L:%f",
		currentKline.CloseTime,
		strategyBO.ToStandardId(_id),
		dir.ToString(),
		quantity,
		currentKline.Close,
		stopProfit,
		stopLoss))

	if sendNotify {
		//send slack
	}
}

func ExitOrder(
	strategyBO *core.StrategyBO,
	_id string,
) {
	order, exists := orderMap[strategyBO.ToStandardId(_id)]
	if !exists {
		return
	}

	order.Exit(currentKline.Close, currentKline.CloseTime)
}

func CancelOrder(
	strategyBO *core.StrategyBO,
	_id string,
	sendNotify bool,
) {
	order, exists := orderMap[strategyBO.ToStandardId(_id)]
	if !exists {
		return
	}

	order.Cancel()

	slog.Info(fmt.Sprintf("[%s][%s] cancelled",
		currentKline.CloseTime,
		strategyBO.ToStandardId(_id),
	))

	if sendNotify {
		//send slack
	}
}

func GetOrderStatus(strategyBO *core.StrategyBO, id string) core.OrderStatus {
	order, exists := orderMap[strategyBO.ToStandardId(id)]

	if exists {
		return order.GetStatus()
	}
	return core.ORDER_UNKWN
}

func orderPut(id string, newOrder *core.OrderBO) bool {
	if order, exists := orderMap[id]; exists {
		if order.GetStatus() == core.ORDER_OPEN {
			orderMap[id] = newOrder
			return true
		} else {
			return false
		}
	}

	orderMap[id] = newOrder
	return true
}

func OutputOrdersResult() {
	f, _ := os.OpenFile(time.Now().Format("20060102150405")+"_report.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	var orderKeys []string
	for k := range orderMap {
		orderKeys = append(orderKeys, k)
	}
	sort.Strings(orderKeys)

	for _, k := range orderKeys {
		f.Write([]byte(fmt.Sprintf("%s\n", orderMap[k].ToCsv())))
	}
}
