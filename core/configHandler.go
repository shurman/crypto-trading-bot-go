package core

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
	Symbols           []string `mapstructure:"symbols"`
	Interval          string   `mapstructure:"interval"`
	InitialFund       float64  `mapstructure:"initialFund"`
	SingleRiskRatio   float64  `mapstructure:"singleRiskRatio"`
	ProfitLossRatio   float64  `mapstructure:"ProfitLossRatio"`
	EnableAccumulated bool     `mapstructure:"enableAccumulated"`
	FeeMakerRate      float64  `mapstructure:"feeMakerRate"`
	FeeTakerRate      float64  `mapstructure:"feeTakerRate"`

	Mode        string             `mapstructure:"mode"`
	Indicator   IndicatorConfigs   `mapstructure:"indicator"`
	Backtesting BacktestingConfigs `mapstructure:"backtesting"`
}

type IndicatorConfigs struct {
	StartFromKlines int `mapstructure:"startFromKlines"`
}

type BacktestingConfigs struct {
	Export       BacktestingExportCsvConfigs `mapstructure:"export"`
	Download     BacktestingDownloadConfigs  `mapstructure:"download"`
	OverallChart bool                        `mapstructure:"overallChart"`
}

type BacktestingExportCsvConfigs struct {
	Orders  bool `mapstructure:"orders"`
	Reports bool `mapstructure:"reports"`
	Chart   bool `mapstructure:"chart"`
}

type BacktestingDownloadConfigs struct {
	Enable           bool  `mapstructure:"enable"`
	StartTime        int64 `mapstructure:"startTime"`
	LimitPerDownload int   `mapstructure:"limitPerDownload"`
}

type BinanceConfigs struct {
	Apikey    string `mapstructure:"apikey"`
	Apisecret string `mapstructure:"apisecret"`
}

type SlackConfigs struct {
	Enable  bool   `mapstructure:"enable"`
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
