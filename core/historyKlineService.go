// historyKlineService
package core

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

var (
	futuresClient *futures.Client
)

func FutureClientInit() {
	futuresClient = binance.NewFuturesClient(bKey, bSecret)
}

func GetHistoryKline() {
	Logger.Info("[GetHistoryKline] Start")

	historyKlines, _ := futuresClient.NewKlinesService().Symbol(tSymbol).Limit(tHistoryLimit + 1).Interval(tInterval).Do(context.Background())

	for _, fKline := range historyKlines[:len(historyKlines)-1] { //remove last one because of not closed
		Logger.Debug(fmt.Sprintf("%+v", fKline))

		hKline := kConvertKline(fKline)
		recordNewKline(&hKline)
	}
}
