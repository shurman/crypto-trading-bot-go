// historyKlineService
package service

import (
	"context"
	"crypto-trading-bot-go/core"
	"fmt"
	"io"
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
	futuresClient = binance.NewFuturesClient(core.Config.Binance.Apikey, core.Config.Binance.Apisecret)
}

func LoadHistoryKline() {
	for _, symbol := range core.Config.Trading.Symbols {
		historyKlines, _ := futuresClient.NewKlinesService().
			Symbol(symbol).
			Limit(max(min(core.Config.Trading.Indicator.StartFromKlines+1, 1500), 1)).
			Interval(core.Config.Trading.Interval).
			Do(context.Background())

		for _, fKline := range historyKlines[:len(historyKlines)-1] { //remove last one because of not closed
			hKline := core.ConvertKlineFromFuturesKline(fKline)

			Logger.Debug("Loaded " + fmt.Sprintf("%+v", hKline))
			recordNewKline(symbol, &hKline)
		}

		Logger.Info(symbol + " Loaded")
	}

	Logger.Info("[LoadHistoryKline] History klines loaded")
}

func DownloadRawHistoryKline(symbol string, interval string, startTime int64, limit int64) {
	f, _ := os.OpenFile(symbol+"_"+interval+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	for {
		Logger.Info(fmt.Sprintf("Downloading %s %d klines from %d", symbol, limit, startTime))
		response, error := http.Get(fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s&startTime=%d&limit=%d", symbol, interval, startTime, limit))

		if error != nil {
			panic(error)
		}

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		json := loadJson(body)

		count := len(json.MustArray())
		if count <= 1 {
			break
		} else if count < 1500 {
			count -= 1
		}

		for i := 0; i < count; i++ {
			item := json.GetIndex(i)

			if len(item.MustArray()) < 11 {
				panic("Raw Kline format error")
			}

			kline := &futures.Kline{
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
			}

			f.Write([]byte(fmt.Sprintf("[%d,%s,%s,%s,%s,%d],\n",
				kline.OpenTime,
				kline.Open,
				kline.High,
				kline.Low,
				kline.Close,
				kline.CloseTime)))

			startTime = kline.CloseTime
		}

		time.Sleep(1 * time.Second)
	}
}

func LoadRawHistoryKline(symbol string, interval string) {
	data, err := os.ReadFile(symbol + "_" + interval + ".txt")

	if err != nil {
		panic(err)
	}

	rawKlines := "[" + strings.ReplaceAll(string(data), "\n", "") + "[]]"
	json := loadJson([]byte(rawKlines))

	count := len(json.MustArray()) - 1

	for i := 0; i < count; i++ {
		item := json.GetIndex(i)

		if len(item.MustArray()) < 6 {
			panic("Kline format error")
		}

		kline := &core.Kline{
			StartTime: time.Unix(item.GetIndex(0).MustInt64()/1000, 0),
			Open:      item.GetIndex(1).MustFloat64(),
			High:      item.GetIndex(2).MustFloat64(),
			Low:       item.GetIndex(3).MustFloat64(),
			Close:     item.GetIndex(4).MustFloat64(),
			CloseTime: time.Unix(item.GetIndex(5).MustInt64()/1000, 0),
		}

		recordNewKline(symbol, kline)
	}
}

func loadJson(rawData []byte) *simplejson.Json {
	json, err := simplejson.NewJson(rawData)

	if err != nil {
		panic(err)
	}

	return json
}
