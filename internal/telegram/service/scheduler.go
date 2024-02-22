package service

//func StartScheduler(b *api.Bot) {
//	c := cron.New(cron.WithSeconds())
//	for _, customer := range b.FetchAllCustomers() {
//		customer := customer // capture range variable
//		_, err := c.AddFunc("@every "+customer.SearchInterval.String(), func() {
//			b.ExecuteSearchQueries(&customer)
//		})
//		if err != nil {
//			log.Println("Ошибка при добавлении задачи в планировщик:", err)
//		}
//	}
//	c.Start()
//}
