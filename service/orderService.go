package service

import (
	"crypto-trading-bot-go/core"
	"fmt"
	"log/slog"
)

var (
	orderMap     = make(map[string]*core.OrderBO)
	currentKline *core.Kline
)

func SetCurrentKline(k *core.Kline) {
	currentKline = k
}

func CheckOrderFilled() {

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
		return
	}

	slog.Info(fmt.Sprintf("[%s][%s] %s %f@%f P:%f L:%f",
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

func CancelOrder(
	strategyBO *core.StrategyBO,
	_id string,
	sendNotify bool,
) {
	order, exists := orderMap[strategyBO.ToStandardId(_id)]
	if !exists {
		return
	}

	order.SetStatus(core.ORDER_CANCEL)

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

func orderPut(id string, order *core.OrderBO) bool {
	if order, exists := orderMap[id]; exists {
		if order.GetStatus() == core.ORDER_OPEN {
			orderMap[id] = order
			return true
		} else {
			return false
		}
	}

	orderMap[id] = order
	return true
}
