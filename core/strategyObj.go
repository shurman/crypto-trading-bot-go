// strategyObj
package core

import (
	"fmt"
	"log/slog"
	"time"
)

type StrategyObj struct {
	Name      string
	Execute   func(StrategyObj)
	NextKline chan Kline
}

func (obj *StrategyObj) CreateOrder(
	kline *Kline,
	id string,
	dir OrderDirection,
	quantity float64,
	entry float64,
	stopProfit float64,
	stopLoss float64,
) {
	slog.Info(fmt.Sprintf("[%s][%s-%s] %s %f@%f P:%f L:%f",
		time.Unix(kline.StartTime/1000, 0),
		obj.Name,
		id,
		dir.toString(),
		quantity,
		entry,
		stopProfit,
		stopLoss))

	if kline.IsNew {
		//send slack
	}
}

func (obj *StrategyObj) CancelOrder(kline *Kline, id string) {
	slog.Info(fmt.Sprintf("[%s][%s][%s] cancelled",
		time.Unix(kline.StartTime/1000, 0),
		obj.Name,
		id))

	if kline.IsNew {
		//send slack
	}
}
