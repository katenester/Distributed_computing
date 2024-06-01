package handle

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/katenester/Distributed_computing/internal/model"
	"github.com/katenester/Distributed_computing/internal/orchestrator/ParseExpression"
	"log"
	"net/http"
	"strconv"
)

type Handler interface {
	Register(router *httprouter.Router)
}
type handler struct {
}

func NewHandler() Handler {
	return &handler{}
}

const (
	addExpressionURL = "/api/v1/calculate"
	expressions      = "/api/v1/expressions"
	expression       = "/api/v1/expressions/:id"
	// Общение с агентом
	task = "/internal/task"
)

// Register - регистрация обработчиков handler
func (h *handler) Register(router *httprouter.Router) {
	// регистрируем пути
	router.POST(addExpressionURL, h.AddExpression)
	router.GET(expressions, h.GetListExpressions)
	router.GET(expression, h.GetExpression)
	router.GET(task, h.GiveTask)
	router.POST(task, h.GetResultTask)
}

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
			Status bool    `json:"status"`
			Result float64 `json:"result"`
		} `json:"expressions"`
	}{
		Expressions: make([]struct {
			ID     int     `json:"id"`
			Status bool    `json:"status"`
			Result float64 `json:"result"`
		}, len(list)),
	}
	for i, v := range list {
		var err error
		resp.Expressions[i].Result, err = strconv.ParseFloat(v.Expr, 64)
		resp.Expressions[i].ID = v.Id
		// true - выражение посчитано, т.е. в стеке осталось одно число по алгоритму обратной польской записи, иначе false
		resp.Expressions[i].Status = err == nil
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
	// Декодируем в json
	resp := map[string]struct {
		ID     int     `json:"id"`
		Status bool    `json:"status"`
		Result float64 `json:"result"`
	}{
		"expression": {
			ID:     exp.Id,
			Status: err == nil,
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

// GiveTask - Передача задачи агенту
func (h *handler) GiveTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Ищем задачу для агента
	if task, isTask := storage.FindTask(); isTask {
		resp := map[string]model.Task{
			"task": task,
		}
		log.Println("Сервер: Передаю задачу агенту", resp)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	// Задачи нет
	w.WriteHeader(http.StatusNotFound)
}

// GetResultTask - Получение результата вычисления от агента
func (h *handler) GetResultTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Получаем id и result
	data := struct {
		Id     int     `json:"id"`
		Result float64 `json:"result"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	log.Println("Сервер: Получение результата вычисления от агента", data)
	// Проверяем таску в хранилке , и если такая есть => меняем результат
	if storage.FindAndReplace(data.Id, data.Result) {
		w.WriteHeader(http.StatusCreated)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
