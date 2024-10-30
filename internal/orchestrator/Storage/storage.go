package Storage

import (
	"errors"
	"github.com/katenester/Distributed_computing/internal/config"
	"github.com/katenester/Distributed_computing/internal/model"
	"sync"
	"time"
)

type Storage struct {
	// Выражения в польской записи
	expressions []Expression
	mx          *sync.RWMutex
}

func New() *Storage {
	return &Storage{
		expressions: make([]Expression, 0),
		mx:          &sync.RWMutex{},
	}
}

// GetExpression -Получение выражения
func (s *Storage) GetExpression(idExpressions int) (Expression, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, v := range s.expressions {
		if v.Id == idExpressions {
			return v, nil
		}
	}
	return Expression{}, errors.New("выражения не найдено")
}

// SaveExpression -Добавление выражение в хранилку. Возвращает его id
func (s *Storage) SaveExpression(expression string) int {
	s.mx.Lock()
	defer s.mx.Unlock()
	id := int(time.Now().UnixNano())
	s.expressions = append(s.expressions, Expression{Id: id, Expr: expression, OriginalExpr: expression, Tasks: make([]model.Task, 0), mx: &sync.RWMutex{}, Err: nil})
	// Генерируем таски для задачи.
	s.expressions[len(s.expressions)-1].GenerateTask()
	// Возвращаем последнюю добавленную запись
	return id
}

// GetAllExpression - получение списка всех выражений
func (s *Storage) GetAllExpression() []Expression {
	return s.expressions
}

// FindAndReplace - делает замены решенных подзадач
func (s *Storage) FindAndReplace(id int32, result float64, err error) bool {
	for i := 0; i < len(s.expressions); i++ {
		if s.expressions[i].IsNotSolved() && s.expressions[i].FindAndReplace(id, result, err) {
			return true
		}
	}
	return false
}

// FindTask - ищет задачу для агента
func (s *Storage) FindTask() (model.Task, bool) {
	for i := 0; i < len(s.expressions); i++ {
		// Если задача не решенная
		if s.expressions[i].IsNotSolved() {
			// Выбираем таску
			for j := 0; j < len(s.expressions[i].Tasks); j++ {
				if s.expressions[i].Tasks[j].LastAccess.IsZero() || time.Now().After(s.expressions[i].Tasks[j].LastAccess) {
					// Ставим дедлайн (время операции + издержки)
					s.expressions[i].Tasks[j].LastAccess = time.Now().Add(config.DEADLINE)
					// Отправляем задачу агенту
					return model.Task{
						Id:         s.expressions[i].Tasks[j].Id,
						Arg1:       s.expressions[i].Tasks[j].Arg1,
						Arg2:       s.expressions[i].Tasks[j].Arg2,
						Operation:  s.expressions[i].Tasks[j].Operation,
						LastAccess: s.expressions[i].Tasks[j].LastAccess,
					}, true
				}
			}
		}
	}
	return model.Task{}, false
}
