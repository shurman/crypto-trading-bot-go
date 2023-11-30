// historyKlineService
package service

import (
	"context"
	"crypto-trading-bot-go/core"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"

	"github.com/bitly/go-simplejson"
)

var (
	futuresClient *futures.Client
)

func init() {
	futuresClient = binance.NewFuturesClient(Config.Binance.Apikey, Config.Binance.Apisecret)
}

func LoadHistoryKline() {
	historyKlines, _ := futuresClient.NewKlinesService().
		Symbol(Config.Trading.Symbol).
		Limit(max(min(Config.Trading.Indicator.StartFromKlines+1, 1500), 1)).
		Interval(Config.Trading.Interval).
		Do(context.Background())

	for _, fKline := range historyKlines[:len(historyKlines)-1] { //remove last one because of not closed
		hKline := core.ConvertKlineFromFuturesKline(fKline)

		Logger.Debug("Loaded " + fmt.Sprintf("%+v", hKline))
		recordNewKline(&hKline)
	}

	Logger.Info("[LoadHistoryKline] History klines loaded.")
}

func DownloadRawHistoryKline(symbol string, interval string, startTime int64, limit int64) {
	f, _ := os.OpenFile(symbol+"_"+interval+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	for {
		slog.Info(fmt.Sprintf("Getting %d", startTime))
		response, error := http.Get(fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s&startTime=%d&limit=%d", symbol, interval, startTime, limit))

		if error != nil {
			panic(error)
		}

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		bodyStr := string(body)
		bodyStr = strings.ReplaceAll(bodyStr[1:len(bodyStr)-1], "],[", "],\n[")
		if bodyStr == "" {
			break
		}
		f.Write([]byte(bodyStr + ",\n"))

		time.Sleep(5 * time.Second)
		startTime = time.Unix(startTime/1000, 0).Add(time.Minute * time.Duration(15*limit)).UnixMilli()
	}
}

func LoadRawHistoryKline(symbol string, interval string) {
	data, err := os.ReadFile(symbol + "_" + interval + ".txt")

	if err != nil {
		panic(err)
	}

	rawKlines := "[" + strings.ReplaceAll(string(data), "\n", "") + "[]]"

	json, err := simplejson.NewJson([]byte(rawKlines))

	if err != nil {
		panic(err)
	}

	num := len(json.MustArray()) - 1

	for i := 0; i < num; i++ {
		item := json.GetIndex(i)

		if len(item.MustArray()) < 11 {
			panic("Kline format error")
		}

		kline := core.ConvertKlineFromFuturesKline(&futures.Kline{
			OpenTime:                 item.GetIndex(0).MustInt64(),
			Open:                     item.GetIndex(1).MustString(),
			High:                     item.GetIndex(2).MustString(),
			Low:                      item.GetIndex(3).MustString(),
			Close:                    item.GetIndex(4).MustString(),
			Volume:                   item.GetIndex(5).MustString(),
			CloseTime:                item.GetIndex(6).MustInt64(),
			QuoteAssetVolume:         item.GetIndex(7).MustString(),
			TradeNum:                 item.GetIndex(8).MustInt64(),
			TakerBuyBaseAssetVolume:  item.GetIndex(9).MustString(),
			TakerBuyQuoteAssetVolume: item.GetIndex(10).MustString(),
		})

		recordNewKline(&kline)
	}
}
