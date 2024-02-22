package api

import (
	"fmt"
	"time"
)

func StartTicker(t time.Duration) {
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for range ticker.C {
		makeRequest()
	}
}

func makeRequest() {
	fmt.Println("Выполняется запрос к API...", time.Now())
}
