package main

import (
	uw "UpworkLeadgen/internal/upwork/api"
	"fmt"

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
	//bot := tg.NewBot()
	//bot.StartBot()
	uw.NewConnect()
}
