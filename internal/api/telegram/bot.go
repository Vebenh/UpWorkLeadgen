package telegram

import (
	"log"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	Connection *telebot.Bot
}

func StartBot() {
	b, err := telebot.NewBot(telebot.Settings{
		Token:  viper.GetString("telegram.apiKey"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *telebot.Message) {
		b.Send(m.Sender, "Привет! Я твой Telegram бот.")
	})

	log.Printf("Бот запущен")
	b.Start()
}
