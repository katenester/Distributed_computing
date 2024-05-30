package Storage

import (
	"fmt"
	"main.go/internal/model"
	"strings"
	"sync"
)

type Expression struct {
	Id           int // Пока индекс в слайсе. На доработке нужна генерация id (иначе при удалении всё смешается)
	Expr         string
	OriginalExpr string
	// Задачи для этого выражения
	Tasks []model.Task
	mx    *sync.RWMutex
}

func (e *Expression) FindAndReplace(idTask int, result float64) bool {
	var i int
	for i = 0; i < len(e.Tasks); i++ {
		// Если нашли нужную таску => делаем замену
		if idTask == e.Tasks[i].Id {
			oldSubstr := fmt.Sprintf("%f %f %c", e.Tasks[i].Arg1, e.Tasks[i].Arg2, e.Tasks[i].Operation)
			newSubstr := fmt.Sprintf("%f", result)
			// Делаем замену
			e.mx.Lock()
			e.Expr = strings.Replace(e.Expr, oldSubstr, newSubstr, 1)
			// Меняем местами последний элемент с текущим
			e.Tasks[i], e.Tasks[len(e.Tasks)-1] = e.Tasks[len(e.Tasks)-1], e.Tasks[i]
			// Удаляем последний элемент
			e.Tasks = e.Tasks[:len(e.Tasks)-1]
			e.mx.Unlock()
			return true
		}
	}
	return false
}
