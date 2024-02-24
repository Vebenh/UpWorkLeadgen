package service

import (
	"log"

	"github.com/robfig/cron/v3"

	tg "UpworkLeadgen/internal/telegram/api"
)

type Scheduler struct {
	Bot         *tg.Bot
	Cron        *cron.Cron
	CustomerJob map[int64]*cron.EntryID
}

func NewScheduler(bot *tg.Bot) *Scheduler {
	return &Scheduler{
		Bot:         bot,
		Cron:        cron.New(cron.WithSeconds()), // Инициализация Cron здесь
		CustomerJob: make(map[int64]*cron.EntryID),
	}
}

func (s *Scheduler) StartScheduler() {
	customers, err := s.Bot.Db.FetchAllCustomers()
	if err != nil {
		log.Println("Ошибка при получении данных пользователей", err)
	}

	for _, customer := range customers {
		customer := customer // capture range variable

		if customer.SearchInterval == 0 {
			continue
		}
		s.Bot.SendMessage(customer.TelegramID, customer.SearchInterval.String())
		entryID, err := s.Cron.AddFunc("@every "+customer.SearchInterval.String(), func() {
			s.Bot.SendMessage(customer.TelegramID, customer.SearchInterval.String())
			// TODO Реализовать ExecuteSearchQueries
			//uw.ExecuteSearchQueries(&customer.SearchQueries))

		})
		if err != nil {
			log.Printf("Ошибка при добавлении задачи для пользователя %d: %v", customer.TelegramID, err)
			return
		}

		s.CustomerJob[customer.TelegramID] = &entryID
	}
	s.Cron.Start()

}

func (s *Scheduler) UpdateCustomer(m *tg.UpdateTimeMessage) {
	if entryID, exists := s.CustomerJob[m.TelegramID]; exists {
		s.Cron.Remove(*entryID)
	}

	entryID, err := s.Cron.AddFunc("@every "+m.SearchInterval.String(), func() {
		s.Bot.SendMessage(m.TelegramID, m.SearchInterval.String())
		//TODO Реализовать ExecuteSearchQueries
		//uw.ExecuteSearchQueries(&customer.SearchQueries))

	})
	if err != nil {
		log.Printf("Ошибка при добавлении задачи для пользователя %d: %v", m.TelegramID, err)
		return
	}

	s.CustomerJob[m.TelegramID] = &entryID
}
