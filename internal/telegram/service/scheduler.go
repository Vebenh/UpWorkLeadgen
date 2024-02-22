package service

import (
	"log"

	"github.com/robfig/cron/v3"

	"UpworkLeadgen/internal/telegram/api"
)

func StartScheduler(b *api.Bot) {
	c := cron.New(cron.WithSeconds())
	customers, err := b.Db.FetchAllCustomers()
	if err != nil {
		log.Println("Ошибка при получении данных пользователей", err)
	}

	for _, customer := range customers {
		customer := customer // capture range variable
		_, err := c.AddFunc("@every "+customer.SearchInterval.String(), func() {
			b.ExecuteSearchQueries(&customer)
		})
		if err != nil {
			log.Println("Ошибка при добавлении задачи в планировщик:", err)
		}
	}
	c.Start()
}
