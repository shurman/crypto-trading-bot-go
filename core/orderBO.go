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
	fee         float64
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

func (bo *OrderBO) GetDirection() OrderDirection {
	return bo.dir
}

func (bo *OrderBO) GetCreateTime() time.Time {
	return bo.createTime
}

func (bo *OrderBO) GetExitTime() time.Time {
	return bo.exitTime
}

func (bo *OrderBO) GetStopLossPrice() float64 {
	return bo.stopLoss
}
func (bo *OrderBO) SetStopLossPrice(slPrice float64) {
	bo.stopLoss = slPrice
}

func (bo *OrderBO) GetFinalProfit() float64 {
	return bo.finalProfit
}

func (bo *OrderBO) GetFee() float64 {
	return bo.fee
}

func (bo *OrderBO) Fill(_time time.Time, isTaker bool) {
	bo.status = ORDER_ENTRY
	bo.fillTime = _time

	if isTaker {
		bo.fee = bo.entry * bo.quantity * Config.Trading.FeeTakerRate / 100
	} else {
		bo.fee = bo.entry * bo.quantity * Config.Trading.FeeMakerRate / 100
	}
}

func (bo *OrderBO) Exit(_price float64, _time time.Time, isTaker bool) {
	bo.status = ORDER_EXIT
	bo.exitTime = _time

	bo.finalProfit = (_price - bo.entry) * bo.quantity * float64(bo.dir)

	if isTaker {
		bo.fee += _price * bo.quantity * Config.Trading.FeeTakerRate / 100
	} else {
		bo.fee += _price * bo.quantity * Config.Trading.FeeMakerRate / 100
	}
}

func (bo *OrderBO) Cancel() {
	bo.status = ORDER_CANCEL
}

func (bo *OrderBO) ToCsv() string {
	return fmt.Sprintf("%s,%s,%s,%f,%f,%s,%f,%f",
		bo.fillTime,
		bo.id,
		bo.dir.ToString(),
		bo.entry,
		bo.quantity,
		bo.exitTime,
		bo.finalProfit,
		bo.fee)
}

func (bo *OrderBO) IsExistedAndFilled() bool {
	if bo == nil {
		return false
	}

	return bo.status == ORDER_ENTRY || bo.status == ORDER_EXIT
}

func (bo *OrderBO) IsInstantFillAndExit() bool {
	return bo.fillTime == bo.exitTime
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
