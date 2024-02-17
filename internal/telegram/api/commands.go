package api

import (
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

const (
	StartCommand      = "/start"
	SearchCommand     = "/search"
	UpdateTimeCommand = "/update"
	HelpCommand       = "/help"
)

const (
	SearchState     = "search_state"
	UpdateTimeState = "update_state"
)

type UserState struct {
	CurrentCommand string
	CurrentFilter  string
	UpdateTime     *time.Ticker
}

var userStates = make(map[int64]*UserState)

func (b *Bot) StartHendler(m *telebot.Message) {

	if _, exists := userStates[m.Sender.ID]; !exists {
		userStates[m.Sender.ID] = &UserState{}
	}

	btn1 := telebot.ReplyButton{Text: "Кнопка 1"}
	btn2 := telebot.ReplyButton{Text: "Кнопка 2"}
	btn3 := telebot.ReplyButton{Text: "Кнопка 3"}
	btn4 := telebot.ReplyButton{Text: "Кнопка 4"}

	markup := &telebot.ReplyMarkup{}
	markup.ReplyKeyboard = [][]telebot.ReplyButton{
		{btn1, btn2, btn3, btn4},
	}
	markup.ResizeReplyKeyboard = true
	markup.OneTimeKeyboard = false // Скрыть клавиатуру после использования

	_, err := b.Connection.Send(m.Sender, "Привет! Выберите действие:", markup)
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}

func (b *Bot) SearchHendler(m *telebot.Message) {
	userStates[m.Sender.ID].CurrentCommand = SearchState
	_, err := b.Connection.Send(m.Sender, "Пожалуйста, введите текст:")
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}

func (b *Bot) UpdateTimeHendler(m *telebot.Message) {
	userStates[m.Sender.ID].CurrentCommand = UpdateTimeState
	_, err := b.Connection.Send(m.Sender, "Введите временной промежуток для обновлений в минутах:")
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}

func (b *Bot) HelpHendler(m *telebot.Message) {
	b.Connection.Send(m.Sender, "HelpHendler")
}

func (b *Bot) TextHendler(m *telebot.Message) {
	state := getUserState(m.Sender.ID)
	switch state.CurrentCommand {
	case SearchState:
		userStates[m.Sender.ID].CurrentFilter = m.Text
		b.Connection.Send(m.Sender, "Поиск сохранен")
		b.Connection.Send(m.Sender, m.Text)
	case UpdateTimeState:
		userStates[m.Sender.ID].UpdateTime = createTickerFromText(m.Text)
		b.Connection.Send(m.Sender, "Временной интервал установлен")
	default:
		b.Connection.Send(m.Sender, "Неизвестная команда")
	}
}
