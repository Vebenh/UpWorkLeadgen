package api

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"

	"UpworkLeadgen/internal/telegram/repository"
	"gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	TgConnection *telebot.Bot
	Db           *repository.Db
}

func NewBot() *Bot {
	bot := &Bot{}

	connection, err := telebot.NewBot(telebot.Settings{
		Token:  viper.GetString("telegram.apiKey"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal("Error connecting to telegram bot:", err)
	}

	bot.TgConnection = connection
	bot.Db.Connection, err = repository.InitDB()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	bot.registerHandlers()

	return bot
}

func (b *Bot) registerHandlers() {
	b.TgConnection.Handle(StartCommand, b.StartHendler)
	b.TgConnection.Handle(SearchCommand, b.SearchHendler)
	b.TgConnection.Handle(UpdateTimeCommand, b.UpdateTimeHandler)
	b.TgConnection.Handle(HelpCommand, b.HelpHendler)
	b.TgConnection.Handle(telebot.OnText, b.TextHendler)
}

func (b *Bot) StartBot() {
	fmt.Println("The bot is running...")
	b.TgConnection.Start()
}
