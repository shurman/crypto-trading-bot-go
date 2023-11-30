// configService.go
package service

import (
	"github.com/spf13/viper"
)

var (
	Config Configurations
)

type Configurations struct {
	System  SystemConfigs  `mapstructure:"system"`
	Trading TradingConfigs `mapstructure:"trading"`
	Binance BinanceConfigs `mapstructure:"binance"`
	Slack   SlackConfigs   `mapstructure:"slack"`
}

type SystemConfigs struct {
	Loglevel string `mapstructure:"loglevel"`
}

type TradingConfigs struct {
	Symbol   string `mapstructure:"symbol"`
	Interval string `mapstructure:"interval"`
	Mode     string `mapstructure:"mode"`

	Indicator   IndicatorConfigs   `mapstructure:"indicator"`
	Backtesting BacktestingConfigs `mapstructure:"backtesting"`
}

type IndicatorConfigs struct {
	StartFromKlines int `mapstructure:"startFromKlines"`
}

type BacktestingConfigs struct {
	Download BacktestingDownloadConfigs `mapstructure:"download"`
}

type BacktestingDownloadConfigs struct {
	Enable           bool  `mapstructure:"enable"`
	StartTime        int64 `mapstructure:"startTime"`
	LimitPerDownload int64 `mapstructure:"limitPerDownload"`
}

type BinanceConfigs struct {
	Apikey    string `mapstructure:"apikey"`
	Apisecret string `mapstructure:"apisecret"`
}

type SlackConfigs struct {
	Webhook string `mapstructure:"webhook"`
	Channel string `mapstructure:"channel"`
}

func init() {
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
}
