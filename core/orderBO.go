package core

import (
	"fmt"
	"log/slog"
	"time"
)

var (
	orderMap     = make(map[string]*orderBO)
	currentKline *Kline
)

type orderBO struct {
	createTime  time.Time
	id          string
	status      OrderStatus
	dir         OrderDirection
	quantity    float64
	entry       float64
	stopProfit  float64
	stopLoss    float64
	fillTime    int64
	finalProfit float64
}

type OrderStatus string

const (
	UNKWN  OrderStatus = "UNKWN"
	OPEN   OrderStatus = "OPEN"
	ENTRY  OrderStatus = "ENTRY"
	EXIT   OrderStatus = "EXIT"
	CANCEL OrderStatus = "CANCEL"
)

func SetCurrentKline(k *Kline) {
	currentKline = k
}

func CheckOrderFilled() {

}

func toStandardId(strategyBO *StrategyBO, id string) string {
	return strategyBO.GetName() + "-" + id
}

func CreateOrder(
	strategyBO *StrategyBO,
	id string,
	dir OrderDirection,
	quantity float64,
	entry float64,
	stopProfit float64,
	stopLoss float64,
	sendNotify bool,
) {

	newOrder := &orderBO{
		createTime: currentKline.CloseTime,
		id:         toStandardId(strategyBO, id),
		status:     OPEN,
		dir:        dir,
		quantity:   quantity,
		entry:      entry,
		stopProfit: stopProfit,
		stopLoss:   stopLoss,
	}

	if !orderPut(newOrder.id, newOrder) {
		return
	}

	slog.Info(fmt.Sprintf("[%s][%s-%s] %s %f@%f P:%f L:%f",
		currentKline.CloseTime,
		strategyBO.GetName(),
		id,
		dir.toString(),
		quantity,
		entry,
		stopProfit,
		stopLoss))

	if sendNotify {
		//send slack
	}

}

func CancelOrder(
	strategyBO *StrategyBO,
	id string,
	sendNotify bool,
) {
	order, exists := orderMap[toStandardId(strategyBO, id)]
	if !exists {
		return
	}

	order.status = CANCEL

	slog.Info(fmt.Sprintf("[%s][%s][%s] cancelled",
		currentKline.CloseTime,
		strategyBO.GetName(),
		id))

	if sendNotify {
		//send slack
	}
}

func GetOrderStatus(strategyBO *StrategyBO, id string) OrderStatus {
	order, exists := orderMap[toStandardId(strategyBO, id)]

	if exists {
		return order.status
	}
	return UNKWN
}

func orderPut(id string, order *orderBO) bool {
	if order, exists := orderMap[id]; exists {
		if order.status == OPEN {
			orderMap[id] = order
			return true
		} else {
			return false
		}
	}

	orderMap[id] = order
	return true
}
