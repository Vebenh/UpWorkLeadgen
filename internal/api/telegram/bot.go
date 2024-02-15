package telegram

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	Connection *telebot.Bot
}

func NewBot() *Bot {
	bot := &Bot{}

	connection, err := telebot.NewBot(telebot.Settings{
		Token:  viper.GetString("telegram.apiKey"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal("Error connecting to telegram bot", err)
	}

	bot.Connection = connection
	bot.registerHandlers()

	return bot
}

func (b *Bot) registerHandlers() {
	b.Connection.Handle(StartCommand, b.StartHendler)
	b.Connection.Handle(SearchCommand, b.SearchHendler)
	b.Connection.Handle(UpdateCommand, b.UpdateHendler)
	b.Connection.Handle(HelpCommand, b.HelpHendler)
	b.Connection.Handle(telebot.OnText, b.TextHendler)
}

func (b *Bot) StartBot() {
	fmt.Println("The bot is running...")
	b.Connection.Start()
}
