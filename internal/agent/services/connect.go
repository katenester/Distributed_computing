package services

import (
	"bytes"
	"encoding/json"
	"io"
	"main.go/internal/model"
	"net/http"
	"time"
)

// Demon - тянет за ручку сервер
func Demon(url []string) {
	for {
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		// Отправляем запрос с контекстом
		req, _ := http.NewRequest("GET", url[0], nil)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		// Получаем код ответа
		statusCode := resp.StatusCode
		if statusCode == 200 {
			// Парсим таску
			buf, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			var task model.Task
			err = json.Unmarshal(buf, &task)
			if err != nil {
				continue
			}
			// Вычисляем
			RequestTaskBody := struct {
				id     int
				result float64
			}{
				id:     task.Id,
				result: calculation(task),
			}
			jsonBody, err := json.Marshal(RequestTaskBody)
			req, err := http.NewRequest("POST", url[1], bytes.NewBuffer(jsonBody))
			_, err = client.Do(req)
			continue
		}
	}
}
