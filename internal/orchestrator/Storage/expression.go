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
	// Ошибка вычисления
	Err error
	// Задачи для этого выражения
	Tasks []model.Task
	mx    *sync.RWMutex
}

// FindAndReplace - нахождение таски и её замена
func (e *Expression) FindAndReplace(idTask int, result float64, err error) bool {
	log.Println("До замены", e.Expr)
	// Ищем таску
	for i := 0; i < len(e.Tasks); i++ {
		// Если нашли нужную таску => делаем замену
		if idTask == e.Tasks[i].Id {
			if err != nil {
				e.Err = err
				return true
			}
			// Делаем замену
			e.replace(i, result)
			// Удаляем текущую таску
			e.Tasks[i], e.Tasks[len(e.Tasks)-1] = e.Tasks[len(e.Tasks)-1], e.Tasks[i]
			e.Tasks = e.Tasks[:len(e.Tasks)-1]
			return true
		}
	}
	return false
}

// replace - заменяет посчитанный результат
func (e *Expression) replace(i int, result float64) {
	s := strings.Split(e.Expr, " ")
	// Ищем нужную операцию и числа
	for j := 2; j < len(s); j++ {
		// Если нашли => делаем замену
		if s[j-2] == fmt.Sprint(e.Tasks[i].Arg1) && s[j-1] == fmt.Sprint(e.Tasks[i].Arg2) && s[j] == fmt.Sprint(e.Tasks[i].Operation) {
			e.mx.Lock()
			// Убираем числа и ставим результат
			s = append(s[:j-2], append([]string{fmt.Sprint(result)}, s[j+1:]...)...)
			// result находится на позиции j-2
			// Значит два случая таски для нового числа
			// Обрабатываем первый случай таски:float(j-3) float(j-2) operand(j-1)
			if j-3 >= 0 && j-1 < len(s) {
				a, err1 := strconv.ParseFloat(s[j-3], 64)
				b, err2 := strconv.ParseFloat(s[j-2], 64)
				err3 := s[j-1] == "+" || s[j-1] == "-" || s[j-1] == "*" || s[j-1] == "/"
				// Если нужная комбинация
				if err1 == nil && err2 == nil && err3 {
					// Создаем, если возможно нужную таску
					id := int(time.Now().UnixNano())
					e.Tasks = append(e.Tasks, model.Task{Id: id, Arg1: a, Arg2: b, Operation: s[j-1]})
				}
			}
			// Обрабатываем второй случай таски : float(j-2) float(j-1) operand(j)
			if j-2 >= 0 && j < len(s) {
				a, err1 := strconv.ParseFloat(s[j-2], 64)
				b, err2 := strconv.ParseFloat(s[j-1], 64)
				err3 := s[j] == "+" || s[j] == "-" || s[j] == "*" || s[j] == "/"
				// Если нужная комбинация
				if err1 == nil && err2 == nil && err3 {
					// Создаем, если возможно нужную таску
					id := int(time.Now().UnixNano())
					e.Tasks = append(e.Tasks, model.Task{Id: id, Arg1: a, Arg2: b, Operation: s[j]})
				}
			}
			e.Expr = strings.Join(s, " ")
			log.Println("После замены", e.Expr)
			e.mx.Unlock()
		}
	}
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
	return len(strings.Join([]string{e.Expr}, " ")) > 1 && e.Err == nil
}
