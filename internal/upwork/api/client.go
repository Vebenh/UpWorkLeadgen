package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/upwork/golang-upwork-oauth2/api"
	"io"
	"log"
	"net/http"
	"os"
)

const cfgFile = "./config/upwork.json"

func NewConnect() {
	ctx := context.Background() // Инициализация контекста

	// Инициализация клиента API с конфигурацией из файла
	client := api.Setup(api.ReadConfig(cfgFile))

	// Проверка наличия доступного токена доступа
	if !client.HasAccessToken(ctx) {
		// Получение URL для авторизации пользователя
		authURL := client.GetAuthorizationUrl("random-state")
		fmt.Println("Visit the authorization url and provide oauth_verifier for further authorization:", authURL)

		// Ввод кода авторизации пользователем
		fmt.Print("Enter the code from URL: ")
		reader := bufio.NewReader(os.Stdin)
		authCode, _ := reader.ReadString('\n')

		// Обмен кода на токен доступа
		token := client.GetToken(ctx, authCode)
		viper.Set("access_token", token.AccessToken)
		fmt.Printf("Access Token: %s\n", token.AccessToken)
	} else {
		fmt.Println("Already have an access token.")
	}

	token := viper.GetString("access_token")
	GraphQLQuery(token)
}

func GraphQLQuery(token string) {
	query, err := os.ReadFile("./graphql/marketplaceJobPostings.graphql")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	variables := `{
        "marketPlaceJobFilter": {
            "searchExpression_eq": "Golang"
        },
        "searchType": "USER_JOBS_SEARCH",
        "sortAttributes": [
            {
                "field": "RECENCY"
            }
        ]
    }`

	body := fmt.Sprintf(`{"query": %q, "variables": %s}`, string(query), variables)

	req, err := http.NewRequest("POST", "https://api.upwork.com/graphql", bytes.NewBufferString(body))
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		log.Fatalf("Ошибка при десериализации: %v", err)
	}

	fmt.Printf("Ответ от Upwork API:\n%s\n", result)
}
