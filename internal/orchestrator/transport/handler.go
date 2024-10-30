package transport

import (
	"github.com/julienschmidt/httprouter"
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
	// Общение с агентом (замена на grpc)
	//task = "/internal/task"
)

// Register - регистрация обработчиков handler
func (h *handler) Register(router *httprouter.Router) {
	// регистрируем пути
	router.POST(addExpressionURL, h.AddExpression)
	router.GET(expressions, h.GetListExpressions)
	router.GET(expression, h.GetExpression)
	//router.GET(task, h.GiveTask)
	//router.POST(task, h.GetResultTask)
}
