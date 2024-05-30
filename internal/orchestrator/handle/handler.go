package handle

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"main.go/internal/config"
	"main.go/internal/orchestrator/ParseExpression"
	"main.go/internal/orchestrator/Storage"
	"net/http"
	"strconv"
	"time"
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
	// Получаем ссылку на слайс
	copyList := storage.GetAllExpression()
	list := make([]Storage.Expression, len(*copyList), cap(*copyList))
	// Копируем содержимое для параллельной обработки
	copy(list, *copyList)
	// Создаем слайс анонимных структур для хранения данных JSON
	resp := struct {
		Expressions []struct {
			ID     int  `json:"id"`
			Status bool `json:"status"`
			Result int  `json:"result"`
		} `json:"expressions"`
	}{
		Expressions: make([]struct {
			ID     int  `json:"id"`
			Status bool `json:"status"`
			Result int  `json:"result"`
		}, len(list)),
	}
	for i, v := range list {
		result, err := strconv.Atoi(v.Expr)
		resp.Expressions[i].ID = v.Id
		// true - выражение посчитано, т.е. в стеке осталось одно число по алгоритму обратной польской записи, иначе false
		resp.Expressions[i].Status = err == nil
		resp.Expressions[i].Result = result
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
	// Получаем id из пути и конвертируем в int
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Получаем выражение по id
	exp, err := storage.GetExpression(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := strconv.Atoi(exp.Expr)
	// Декодируем в json
	resp := map[string]struct {
		ID     int  `json:"id"`
		Status bool `json:"status"`
		Result int  `json:"result"`
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
	// Пока для проверки работы будем сканировать все задачи
	// Но нужно будет переделать (например две хранилки - решенные задачи и в процессе)
	list := *storage.GetAllExpression()
	// Ищем задачи с приоритетом
	for _, v := range list {
		// Если задача не решенная
		if len(v.Expr) != 1 {
			// Выбираем таску
			for _, taskExpression := range v.Tasks {
				// Если задача не была взята другим агентом (т.е. не запустился таймер)
				// Либо вышел таймаут
				if taskExpression.LastAccess.IsZero() || time.Now().After(taskExpression.LastAccess) {
					// Ставим дедлайн (время операции + издержки)
					taskExpression.LastAccess = time.Now().Add(time.Duration(config.DEADLINE))
					// Отправляем задачу агенту
					// Декодируем в json
					resp := map[string]struct {
						ID         int           `json:"id"`
						Arg1       float64       `json:"arg1"`
						Arg2       float64       `json:"arg2"`
						Operation  byte          `json:"operation"`
						LastAccess time.Duration `json:"operation_time"`
					}{
						"task": {
							ID:         taskExpression.Id,
							Arg1:       taskExpression.Arg1,
							Arg2:       taskExpression.Arg2,
							Operation:  taskExpression.Operation,
							LastAccess: config.GetDuration(string(taskExpression.Operation)),
						},
					}
					w.Header().Set("Content-Type", "application/json")
					err := json.NewEncoder(w).Encode(resp)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusCreated)
					return
				}
			}
		}
	}
	// Задачи нет
	w.WriteHeader(http.StatusNotFound)
}

// GetResultTask - Получение результата вычисления от агента
func (h *handler) GetResultTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Получаем id и result
	data := struct {
		id     int
		result float64
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// Записываем результат в выражение
	// Обновляем таски
	list := *storage.GetAllExpression()
	// Ищем таску
	for _, v := range list {
		// Если задача не решенная и есть нужная таска
		if len(v.Expr) != 1 && v.FindAndReplace(data.id, data.result) {
			w.WriteHeader(http.StatusCreated)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
