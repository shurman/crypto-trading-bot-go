// configService.go
package core

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	Config Configurations
)

type Configurations struct {
	Trading TradingConfigs `mapstructure:"trading"`
	Binance BinanceConfigs `mapstructure:"binance"`
	Slack   SlackConfigs   `mapstructure:"slack"`
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

type SlackConfigs struct {
	Webhook string `mapstructure:"webhook"`
	Channel string `mapstructure:"channel"`
}

func LoadConfigs() {
	Logger.Info("[LoadConfigs] Start")

	reader := viper.New()

	reader.SetConfigName("config")
	reader.AddConfigPath("./")
	reader.AutomaticEnv()

	reader.SetConfigType("yml")

	if err := reader.ReadInConfig(); err != nil {
		panic("Error reading config file, " + err.Error())
	}

	err := reader.Unmarshal(&Config)
	if err != nil {
		panic("Unable to decode into struct, " + err.Error())
	}

	Logger.Debug(fmt.Sprintf("%+v", Config))
}
