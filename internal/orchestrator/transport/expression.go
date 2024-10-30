package transport

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/katenester/Distributed_computing/internal/orchestrator/ParseExpression"
	"log"
	"net/http"
	"strconv"
)

// AddExpression - добавление выражения
func (h *handler) AddExpression(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Сервер: операция добавления выражения")
	// Создаем анонимную структуру для хранения данных JSON
	var req struct {
		Expression string `json:"expression"`
	}
	// Декодирование JSON тела запроса в карту
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Проверка на валидность и перевод в обратную польскую запись
	polishExpression, err := ParseExpression.ToPostfix(req.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// Сохраняем выражение в хранилище и получаем его ID
	id := storage.SaveExpression(polishExpression)
	// Формируем и отправляем JSON-ответ
	resp := struct {
		ID int `json:"id"`
	}{
		ID: id,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	// Разбить сразу на таски для этого выражения
}

// GetListExpressions - Получение списка всех выражений
func (h *handler) GetListExpressions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Сервер: получение списка всех выражений")
	// Получаем ссылку на слайс
	list := storage.GetAllExpression()
	log.Println("Сервер: хранилка", list)
	// Создаем слайс анонимных структур для хранения данных JSON
	resp := struct {
		Expressions []struct {
			ID     int     `json:"id"`
			Status string  `json:"status"`
			Result float64 `json:"result"`
		} `json:"expressions"`
	}{
		Expressions: make([]struct {
			ID     int     `json:"id"`
			Status string  `json:"status"`
			Result float64 `json:"result"`
		}, len(list)),
	}
	for i, v := range list {
		var err error
		resp.Expressions[i].Result, err = strconv.ParseFloat(v.Expr, 64)
		resp.Expressions[i].ID = v.Id
		// Решено, В процессе, Ошибка : деление на ноль
		if v.Err != nil {
			resp.Expressions[i].Status = "Ошибка : деление на ноль"
		} else if err != nil {
			resp.Expressions[i].Status = "В процессе"
		} else {
			resp.Expressions[i].Status = "Решено"
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetExpression - Получение выражения по id
func (h *handler) GetExpression(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("Сервер: Получение выражения по id")
	// Получаем id из пути и конвертируем в int
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Получаем выражение по id
	exp, err := storage.GetExpression(id)
	log.Println("Сервер: хранилка id", exp)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res, err := strconv.ParseFloat(exp.Expr, 64)
	var status string
	if exp.Err != nil {
		status = "Ошибка: деление на ноль"
	} else if err != nil {
		status = "В процессе"
	} else {
		status = "Решено"
	}
	// Декодируем в json
	resp := map[string]struct {
		ID     int     `json:"id"`
		Status string  `json:"status"`
		Result float64 `json:"result"`
	}{
		"expression": {
			ID:     exp.Id,
			Status: status,
			Result: res,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
