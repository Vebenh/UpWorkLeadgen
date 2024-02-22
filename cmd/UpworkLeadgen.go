package main

import (
	tg "UpworkLeadgen/internal/telegram/api"
	"UpworkLeadgen/internal/telegram/service"
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
	entry := make(chan *tg.Bot)

	wg.Add(2)
	go func() {
		bot := tg.NewBot()
		entry <- bot
		bot.StartBot()
		wg.Done()
	}()
	go func() {
		service.StartScheduler(<-entry)
		wg.Done()
	}()
	wg.Wait()
}
