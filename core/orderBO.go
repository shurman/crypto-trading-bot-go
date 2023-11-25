package core

import (
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
	fillTime    int64
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

func (bo *OrderBO) SetStatus(_status OrderStatus) {
	bo.status = _status
}

type OrderStatus string

const (
	ORDER_UNKWN  OrderStatus = "UNKWN"
	ORDER_OPEN   OrderStatus = "OPEN"
	ORDER_ENTRY  OrderStatus = "ENTRY"
	ORDER_EXIT   OrderStatus = "EXIT"
	ORDER_CANCEL OrderStatus = "CANCEL"
)

type OrderDirection int

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
