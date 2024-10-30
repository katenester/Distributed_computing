package model

import (
	"time"
)

// Task - задачи для вычислений агента
type Task struct {
	Id         int32     `json:"id"`         // идентификатор задачи
	Arg1       float64   `json:"x"`          // имя первого аргумента
	Arg2       float64   `json:"y"`          //имя второго аргумента
	Operation  string    `json:"operator"`   //операция
	LastAccess time.Time `json:"LastAccess"` // время истечения срока выполнения задачи
}
