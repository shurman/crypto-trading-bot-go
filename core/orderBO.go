package core

import (
	"fmt"
	"time"
)

type OrderBO struct {
	createTime  time.Time
	id          string
	status      OrderStatus
	dir         OrderDirection
	quantity    float64
	entry       float64
	stopProfit  float64
	stopLoss    float64
	fillTime    time.Time
	exitTime    time.Time
	finalProfit float64
}

func ConstructOrderBO(
	createTime time.Time,
	id string,
	status OrderStatus,
	dir OrderDirection,
	quantity float64,
	entry float64,
	stopProfit float64,
	stopLoss float64,
) *OrderBO {
	return &OrderBO{
		createTime: createTime,
		id:         id,
		status:     ORDER_OPEN,
		dir:        dir,
		quantity:   quantity,
		entry:      entry,
		stopProfit: stopProfit,
		stopLoss:   stopLoss,
	}
}

func (bo *OrderBO) GetId() string {
	return bo.id
}

func (bo *OrderBO) GetStatus() OrderStatus {
	return bo.status
}

func (bo *OrderBO) GetEntryPrice() float64 {
	return bo.entry
}

func (bo *OrderBO) GetStopProfitPrice() float64 {
	return bo.stopProfit
}

func (bo *OrderBO) GetStopLossPrice() float64 {
	return bo.stopLoss
}

func (bo *OrderBO) GetFinalProfit() float64 {
	return bo.finalProfit
}

func (bo *OrderBO) Fill(_time time.Time) {
	bo.status = ORDER_ENTRY
	bo.fillTime = _time
}

func (bo *OrderBO) Exit(_price float64, _time time.Time) {
	bo.status = ORDER_EXIT
	bo.exitTime = _time
	bo.finalProfit = (_price - bo.entry) * bo.quantity * float64(bo.dir)
}

func (bo *OrderBO) Cancel() {
	bo.status = ORDER_CANCEL
}

func (bo *OrderBO) ToCsv() string {
	return fmt.Sprintf("%s,%s,%s,%f,%f,%s,%f",
		bo.fillTime,
		bo.id,
		bo.dir.ToString(),
		bo.entry,
		bo.quantity,
		bo.exitTime,
		bo.finalProfit)
}

type OrderStatus string

const (
	ORDER_UNKWN  OrderStatus = "UNKWN"
	ORDER_OPEN   OrderStatus = "OPEN"
	ORDER_ENTRY  OrderStatus = "ENTRY"
	ORDER_EXIT   OrderStatus = "EXIT"
	ORDER_CANCEL OrderStatus = "CANCEL"
)

type OrderDirection float64

const (
	ORDER_LONG  OrderDirection = 1
	ORDER_SHORT OrderDirection = -1
)

func (dir OrderDirection) ToString() string {
	if dir == ORDER_LONG {
		return "LONG"
	}
	return "SHORT"
}
