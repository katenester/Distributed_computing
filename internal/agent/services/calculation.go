package services

import (
	"errors"
	"github.com/katenester/Distributed_computing/internal/config"
	"github.com/katenester/Distributed_computing/internal/model"
	"time"
)

// Вычисляет выражение "долго"
func calculation(task model.Task) (float64, error) {
	switch string(task.Operation) {
	case "+":
		time.Sleep(config.TIME_ADDITION_MS)
		return task.Arg1 + task.Arg2, nil
	case "-":
		time.Sleep(config.TIME_SUBTRACTION_MS)
		return task.Arg1 - task.Arg2, nil
	case "*":
		time.Sleep(config.TIME_MULTIPLICATIONS_MS)
		return task.Arg1 * task.Arg2, nil
	case "/":
		time.Sleep(config.TIME_DIVISIONS_MS)
		if task.Arg2 == 0 {
			return 0, errors.New("Деление на ноль")
		}
		return task.Arg1 / task.Arg2, nil
	}
	return 0.0, nil
}
