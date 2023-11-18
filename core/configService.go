// configService.go
package core

import (
	"log/slog"

	"github.com/spf13/viper"
)

var (
	tSymbol       string
	tInterval     string
	tHistoryLimit int
	bKey          string
	bSecret       string
)

type Configurations struct {
	Trading TradingConfigs `mapstructure:"trading"`
	Binance BinanceConfigs `mapstructure:"binance"`
}

type TradingConfigs struct {
	Symbol       string `mapstructure:"symbol"`
	Interval     string `mapstructure:"interval"`
	HistoryLimit int    `mapstructure:"historylimit"`
}

type BinanceConfigs struct {
	Apikey    string `mapstructure:"apikey"`
	Apisecret string `mapstructure:"apisecret"`
}

func LoadConfigs() {
	slog.Info("LoadConfigs Start")

	reader := viper.New()

	reader.SetConfigName("config")
	reader.AddConfigPath("./")
	reader.AutomaticEnv()

	reader.SetConfigType("yml")

	if err := reader.ReadInConfig(); err != nil {
		panic("Error reading config file, " + err.Error())
	}

	var config Configurations

	err := reader.Unmarshal(&config)
	if err != nil {
		panic("Unable to decode into struct, " + err.Error())
	}

	tSymbol = config.Trading.Symbol
	tInterval = config.Trading.Interval
	tHistoryLimit = config.Trading.HistoryLimit
	bKey = config.Binance.Apikey
	bSecret = config.Binance.Apisecret
}
