package service

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
)

func GetHistoryOpenInterest(symbol string, period string, endTime int64) *futures.OpenInterestStatistic {

	resultSlice, err := FuturesClient.NewOpenInterestStatisticsService().
		Symbol(symbol).
		Period(period).
		EndTime(endTime).
		Limit(1).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	return resultSlice[0]
}

func GetHistoryLongShortRatio(symbol string, period string, endTime int64) *futures.LongShortRatio {

	resultSlice, err := FuturesClient.NewLongShortRatioService().
		Symbol(symbol).
		Period(period).
		EndTime(endTime).
		Limit(1).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	return resultSlice[0]
}
