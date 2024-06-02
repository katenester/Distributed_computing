package services

import (
	"bytes"
	"encoding/json"
	"github.com/katenester/Distributed_computing/internal/model"
	"io"
	"log"
	"net/http"
	"time"
)

// Demon - тянет за ручку сервер
func Demon(url string, dead chan int, j int) {
	log.Println("Запуск", j, "вычислительной машины")
	for {
		client := http.Client{
			Timeout: 20 * time.Second,
		}
		log.Println("Машина", j, "Отправка get запроса ")
		// Отправляем запрос с контекстом
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Машина", j, "Ошибка запроса get ответа:", err)
			continue
		}
		// Получаем код ответа
		statusCode := resp.StatusCode
		log.Println("Машина", j, "Статус get ответа:", statusCode)
		if statusCode == 200 {
			// Парсим таску
			buf, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Машина", j, "Ошибка чтения :", err)
				continue
			}
			task := make(map[string]model.Task, 1)
			jsonTask := "task"
			err = json.Unmarshal(buf, &task)
			if err != nil {
				log.Println("Машина", j, "Ошибка чтения :", err)
				continue
			}
			log.Println("Машина", j, "Get ответ:", task[jsonTask])
			// Вычисляем
			RequestTaskBody := struct {
				Id     int     `json:"id"`
				Result float64 `json:"result"`
				Err    error   `json:"error"`
			}{}
			RequestTaskBody.Id = task[jsonTask].Id
			RequestTaskBody.Result, RequestTaskBody.Err = calculation(task[jsonTask])
			log.Println("Машина", j, "Результат долгой операции:", RequestTaskBody.Result)
			jsonBody, err := json.Marshal(RequestTaskBody)
			// Отправляем её на сервер
			log.Println("Машина", j, "Отправка POST ответа:", RequestTaskBody)
			log.Println("Машина", j, "Отправка POST ответа:", jsonBody)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err = client.Do(req)
			// Получаем код ответа
			statusCode := resp.StatusCode
			log.Println("Машина", j, "Статус post ответа:", statusCode)
			continue
		}
	}
	dead <- 1
}
