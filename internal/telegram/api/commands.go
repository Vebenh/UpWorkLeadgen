package api

import (
	"gopkg.in/tucnak/telebot.v2"
	"log"
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
}

var userStates = make(map[int64]*UserState)

func (b *Bot) StartHandler(m *telebot.Message) {

	if err := b.EnsureUserExists(m.Sender.ID); err != nil {
		log.Fatal("Error verifying the user's existence:", err)
	}

	_, err := b.TgConnection.Send(m.Sender, "Привет! Выберите действие:")
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}

func (b *Bot) SearchHandler(m *telebot.Message) {
	userStates[m.Sender.ID].CurrentCommand = SearchState
	_, err := b.TgConnection.Send(m.Sender, "Пожалуйста, введите текст:")
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}

func (b *Bot) UpdateTimeHandler(m *telebot.Message) {
	userStates[m.Sender.ID].CurrentCommand = UpdateTimeState
	_, err := b.TgConnection.Send(m.Sender, "Введите временной промежуток для обновлений в минутах:")
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}

func (b *Bot) HelpHandler(m *telebot.Message) {
	b.TgConnection.Send(m.Sender, "HelpHendler")
}

func (b *Bot) TextHendler(m *telebot.Message) {
	state := getUserState(m.Sender.ID)
	switch state.CurrentCommand {
	case SearchState:
		b.TgConnection.Send(m.Sender, "Поиск сохранен")
	case UpdateTimeState:
		b.TgConnection.Send(m.Sender, "Временной интервал установлен")
	default:
		b.TgConnection.Send(m.Sender, "Неизвестная команда")
	}
}

func (b *Bot) EnsureUserExists(telegramID int64) error {

	customer, err := b.Db.GetCustomerByTelegramID(telegramID)
	if err != nil {
		return err
	}

	if customer == nil {
		_, err := b.Db.CreateCustomer(telegramID, 0)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
