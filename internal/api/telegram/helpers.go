package telegram

import (
	"log"
	"strconv"
	"time"
)

func getUserState(userID int64) *UserState {
	if state, exists := userStates[userID]; exists {
		return state
	}
	userStates[userID] = &UserState{}
	return userStates[userID]
}

func createTickerFromText(text string) *time.Ticker {
	minutes, err := strconv.Atoi(text)
	if err != nil {
		log.Fatalf("Ошибка при преобразовании текста в число: %v", err)
	}

	// Преобразование минут в Duration
	duration := time.Duration(minutes) * time.Minute

	// Создание и возврат нового Ticker
	ticker := time.NewTicker(duration)
	return ticker
}
