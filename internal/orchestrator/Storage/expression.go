package Storage

import (
	"fmt"
	"github.com/katenester/Distributed_computing/internal/model"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Expression struct {
	Id           int // Пока индекс в слайсе. На доработке нужна генерация id (иначе при удалении всё смешается)
	Expr         string
	OriginalExpr string
	// Задачи для этого выражения
	Tasks []model.Task
	mx    *sync.RWMutex
}

// FindAndReplace - нахождение таски и её замена
func (e *Expression) FindAndReplace(idTask int, result float64) bool {
	log.Println("До замены", e.Expr)
	for i := 0; i < len(e.Tasks); i++ {
		// Если нашли нужную таску => делаем замену
		if idTask == e.Tasks[i].Id {
			s := strings.Split(e.Expr, " ")
			// Ищем нужную операцию и числа
			for j := 2; j < len(s); j++ {
				// Если нашли => делаем замену
				if s[j-2] == fmt.Sprint(e.Tasks[i].Arg1) && s[j-1] == fmt.Sprint(e.Tasks[i].Arg2) && s[j] == fmt.Sprint(e.Tasks[i].Operation) {
					e.mx.Lock()
					// Убираем числа и ставим результат
					s = append(s[:j-2], append([]string{fmt.Sprint(result)}, s[j+1:]...)...)
					// Удаляем текущую таску
					e.Tasks[i], e.Tasks[len(e.Tasks)-1] = e.Tasks[len(e.Tasks)-1], e.Tasks[i]
					e.Tasks = e.Tasks[:len(e.Tasks)-1]
					if j-3 >= 0 {
						// Создаем, если возможно нужную таску
						// P/S/ её можно создать только если после осталось что-то
						// По польской нотации -> там операнд
						id := int(time.Now().UnixNano())
						a, _ := strconv.ParseFloat(s[j-3], 64)
						b, _ := strconv.ParseFloat(s[j-2], 64)
						e.Tasks = append(e.Tasks, model.Task{Id: id, Arg1: a, Arg2: b, Operation: s[j-1]})
					}
					e.Expr = strings.Join(s, " ")
					log.Println("После замены", e.Expr)
					e.mx.Unlock()
					return true
				}
			}
		}
	}
	return false
}

// GenerateTask - генерация новой задачи на таски агенту
func (e *Expression) GenerateTask() {
	e.mx.Lock()
	// Делим на числа и знаки
	s := strings.Split(e.Expr, " ")
	for i := 2; i < len(s); i++ {
		// Проверяем, можно ли сделать таску
		a, err2 := strconv.ParseFloat(s[i-2], 64)
		b, err1 := strconv.ParseFloat(s[i-1], 64)
		if err1 == nil && err2 == nil && (s[i] == "+" || s[i] == "-" || s[i] == "*" || s[i] == "/") {
			id := int(time.Now().UnixNano())

			e.Tasks = append(e.Tasks, model.Task{Id: id, Arg1: a, Arg2: b, Operation: s[i]})
		}
	}
	e.mx.Unlock()
}

// IsNotSolved - Проверяет , не решена ли задача
func (e *Expression) IsNotSolved() bool {
	return len(strings.Join([]string{e.Expr}, " ")) > 1
}
