package main

import (
	tg "UpworkLeadgen/internal/telegram/api"
	"fmt"
	"sync"

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
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		bot := tg.NewBot()
		bot.StartBot()
		wg.Done()
	}()
	go func() {
		fmt.Println("awserdtgfhj")
		wg.Done()
	}()
	wg.Wait()

	//uw.NewConnect()
}
