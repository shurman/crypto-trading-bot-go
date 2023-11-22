package core

var (
	orderSlice []*OrderBO
)

type OrderBO struct {
	createTime  int64
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
	OPEN   OrderStatus = "OPEN"
	ENTRY  OrderStatus = "ENTRY"
	EXIT   OrderStatus = "EXIT"
	CANCEL OrderStatus = "CANCEL"
)
