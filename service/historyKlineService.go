// historyKlineService
package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
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
		Limit(Config.Trading.HistoryLimit + 1).
		Interval(Config.Trading.Interval).
		Do(context.Background())

	for _, fKline := range historyKlines[:len(historyKlines)-1] { //remove last one because of not closed
		hKline := kConvertKline(fKline)

		Logger.Debug("Loaded " + fmt.Sprintf("%+v", hKline))
		recordNewKline(&hKline)
	}

	Logger.Info("[LoadHistoryKline] History klines loaded.")
}

func DownloadRawHistoryKline(symbol string, interval string, startTime int64, limit int64) {
	f, _ := os.Create(symbol + "_" + interval + ".txt")
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
