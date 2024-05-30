package Storage

import (
	"errors"
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
	s.expressions = append(s.expressions, Expression{Id: id, Expr: expression, OriginalExpr: expression})
	// Возвращаем последнюю добавленную запись
	return id
}

// GetAllExpression - получение списка всех выражений
func (s *Storage) GetAllExpression() *[]Expression {
	return &s.expressions
}
