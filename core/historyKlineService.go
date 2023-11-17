// historyKlineService
package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

var (
	futuresClient *futures.Client
)

func FutureClientInit() {
	futuresClient = binance.NewFuturesClient(aKey, aSecret)
}

func GetHistoryKline() {
	historyKlines, _ := futuresClient.NewKlinesService().Symbol(sSymbol).Limit(9).Interval(sInterval).Do(context.Background())

	for _, hKline := range historyKlines {
		slog.Info(fmt.Sprintf("%+v", hKline))
	}

}
