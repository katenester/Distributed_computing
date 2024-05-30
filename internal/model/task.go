package model

import (
	"time"
)

// Task - задачи для вычислений агента
type Task struct {
	Id         int           `json:"id"`             // идентификатор задачи
	Arg1       float64       `json:"arg1"`           // имя первого аргумента
	Arg2       float64       `json:"arg2"`           //имя второго аргумента
	Operation  byte          `json:"operation"`      //операция
	LastAccess time.Duration `json:"operation_time"` // время истечения срока выполнения задачи
}
