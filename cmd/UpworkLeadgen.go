package main

import (
	tg "UpworkLeadgen/internal/telegram/api"
	"UpworkLeadgen/internal/telegram/service"
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
	botChan := make(chan *tg.Bot)
	defer close(botChan)

	go func() {
		bot := <-botChan
		scheduler := service.NewScheduler(bot)
		scheduler.StartScheduler()
		fmt.Println("Горутина шедулера")
		for {
			fmt.Println("1")
			updateTimeMessage := <-bot.UpdateTimeChannel
			fmt.Println("updateTimeMessage дошел")
			scheduler.UpdateCustomer(updateTimeMessage)
			fmt.Println("2")
		}
	}()
	go func() {
		// TODO Реализовать http хендлер для OAuth 2.0
		uw.NewConnect()
	}()

	bot := tg.NewBot()
	botChan <- bot
	bot.StartBot()
}
