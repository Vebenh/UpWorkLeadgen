package main

import (
	"fmt"

	_ "UpworkLeadgen/internal/api"
	"UpworkLeadgen/internal/api/telegram"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func main() {
	bot := telegram.NewBot()
	bot.StartBot()
}
