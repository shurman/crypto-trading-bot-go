// configService.go
package core

import (
	"log/slog"

	"github.com/spf13/viper"
)

var (
	sSymbol   string
	sInterval string
	aKey      string
	aSecret   string
)

type Configurations struct {
	Trading TradingConfigs `mapstructure:"trading"`
	Binance BinanceConfigs `mapstructure:"binance"`
}

type TradingConfigs struct {
	Symbol   string `mapstructure:"symbol"`
	Interval string `mapstructure:"interval"`
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

	sSymbol = config.Trading.Symbol
	sInterval = config.Trading.Interval
	aKey = config.Binance.Apikey
	aSecret = config.Binance.Apisecret
}
