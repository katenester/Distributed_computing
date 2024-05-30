package services

import (
	"main.go/internal/config"
	"main.go/internal/model"
	"time"
)

// Вычисляет выражение "долго"
func calculation(task model.Task) float64 {
	switch string(task.Operation) {
	case "+":
		time.Sleep(config.TIME_ADDITION_MS)
		return task.Arg1 + task.Arg2
	case "-":
		time.Sleep(config.TIME_SUBTRACTION_MS)
		return task.Arg1 - task.Arg2
	case "*":
		time.Sleep(config.TIME_MULTIPLICATIONS_MS)
		return task.Arg1 * task.Arg2
	case "/":
		time.Sleep(config.TIME_DIVISIONS_MS)
		return task.Arg1 / task.Arg2
	}
	return 0.0
}
