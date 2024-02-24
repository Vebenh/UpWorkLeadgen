package api

import (
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
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

	if err := b.Db.EnsureCustomerExists(m.Sender.ID); err != nil {
		log.Fatal("Error verifying the user's existence:", err)
	}
	if _, exists := userStates[m.Sender.ID]; !exists {
		userStates[m.Sender.ID] = &UserState{}
	}

	markup := &telebot.ReplyMarkup{}
	markup.ResizeReplyKeyboard = true
	markup.OneTimeKeyboard = true

	btnSearchProjects := markup.Data("Поисковый запрос", SearchCommand)
	btnSettings := markup.Data("Время обновления", UpdateTimeCommand)

	markup.Inline(markup.Row(btnSearchProjects, btnSettings))

	_, err := b.TgConnection.Send(m.Sender, "Привет! Я твой помощник по поиску работы на Upwork. Выбери действие:", markup)
	if err != nil {
		log.Println("Error sending the message:", err)
	}
}

func (b *Bot) CallbackHandler(c *telebot.Callback) {
	command := strings.TrimSpace(c.Data)

	switch command {
	case UpdateTimeCommand:
		userStates[c.Sender.ID].CurrentCommand = UpdateTimeState
		_, err := b.TgConnection.Send(c.Sender, "Пожалуйста, введите время в минутах:")
		if err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
		}
	case SearchCommand:
		userStates[c.Sender.ID].CurrentCommand = SearchState
		_, err := b.TgConnection.Send(c.Sender, "Пожалуйста, введите поисковый запрос:")
		if err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
		}
	default:
		b.TgConnection.Send(c.Sender, "Неизвестная команда")
	}

	if err := b.TgConnection.Respond(c, &telebot.CallbackResponse{}); err != nil {
		log.Println("Error responding to callback:", err)
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

func (b *Bot) TextHandler(m *telebot.Message) {
	state := getUserState(m.Sender.ID)
	switch state.CurrentCommand {
	case SearchState:
		err := b.Db.AddSearchQuery(m.Sender.ID, m.Text)
		if err != nil {
			log.Println("Ошибка при добавлении поискового запроса:", err)
			return
		}
		b.TgConnection.Send(m.Sender, "Поисковый запрос сохранен")
	case UpdateTimeState:
		duration, err := createDurationFromText(m.Text)
		if err != nil {
			log.Println("Ошибка преобразования времени:", err)
			return
		}
		err = b.Db.SetUpdateTime(m.Sender.ID, duration)
		if err != nil {
			log.Println("Ошибка обновления времени", err)
			return
		}
		b.UpdateTimeChannel <- &UpdateTimeMessage{TelegramID: m.Sender.ID, SearchInterval: duration}
		b.TgConnection.Send(m.Sender, "Временной интервал установлен")
	default:
		b.TgConnection.Send(m.Sender, "Неизвестная команда")
	}
}

func (b *Bot) SendMessage(telegramID int64, message string) {
	recipient := &telebot.User{ID: telegramID}

	_, err := b.TgConnection.Send(recipient, message, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown, // Используйте ModeHTML или ModeMarkdown для форматирования
	})
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v\n", err)
	}
}
