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
	TgConnection      *telebot.Bot
	Db                *repository.Db
	UpdateTimeChannel chan *UpdateTimeMessage
}

type UpdateTimeMessage struct {
	TelegramID     int64
	SearchInterval time.Duration
}

func NewBot() *Bot {
	bot := &Bot{UpdateTimeChannel: make(chan *UpdateTimeMessage)}

	tgConnection, err := telebot.NewBot(telebot.Settings{
		Token:  viper.GetString("telegram.apiKey"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal("Error connecting to telegram bot:", err)
	}

	Db, err := repository.InitDB()

	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	bot.TgConnection = tgConnection
	bot.Db = Db

	bot.registerHandlers()

	commands := []telebot.Command{
		{Text: "/start", Description: "Запустить бота"},
		{Text: "/help", Description: "Показать помощь"},
	}

	err = bot.TgConnection.SetCommands(commands)
	if err != nil {
		log.Fatal(err)
	}

	return bot
}

func (b *Bot) registerHandlers() {
	b.TgConnection.Handle(StartCommand, b.StartHandler)
	b.TgConnection.Handle(SearchCommand, b.SearchHandler)
	b.TgConnection.Handle(UpdateTimeCommand, b.UpdateTimeHandler)
	b.TgConnection.Handle(HelpCommand, b.HelpHandler)
	b.TgConnection.Handle(telebot.OnText, b.TextHandler)
	b.TgConnection.Handle(telebot.OnCallback, b.CallbackHandler)
}

func (b *Bot) StartBot() {
	fmt.Println("The bot is running...")
	b.TgConnection.Start()
}
