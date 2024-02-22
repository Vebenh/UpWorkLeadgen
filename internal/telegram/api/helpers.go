package api

import (
	"fmt"
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

func createDurationFromText(text string) (time.Duration, error) {
	minutes, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("ошибка при преобразовании текста в число: %w", err)
	}

	duration := time.Duration(minutes) * time.Minute
	return duration, nil
}
