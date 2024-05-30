package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

/*Время выполнения операций задается переменными среды в милисекундах

TIME_ADDITION_MS - время выполнения операции сложения в милисекундах
TIME_SUBTRACTION_MS - время выполнения операции вычитания в милисекундах
TIME_MULTIPLICATIONS_MS - время выполнения операции умножения в милисекундах
TIME_DIVISIONS_MS - время выполнения операции деления в милисекундах*/

var TIME_ADDITION_MS, TIME_SUBTRACTION_MS, TIME_MULTIPLICATIONS_MS, TIME_DIVISIONS_MS, DEADLINE time.Duration

func init() {
	TIME_ADDITION_MS = getEnvAsInt("TIME_ADDITION_MS", 2000)
	TIME_SUBTRACTION_MS = getEnvAsInt("TIME_SUBTRACTION_MS", 2000)
	TIME_MULTIPLICATIONS_MS = getEnvAsInt("TIME_MULTIPLICATIONS_MS", 3000)
	TIME_DIVISIONS_MS = getEnvAsInt("TIME_DIVISIONS_MS", 3000)
	DEADLINE = TIME_DIVISIONS_MS * 5
}

func getEnvAsInt(name string, defaultValue int) time.Duration {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return time.Duration(defaultValue)
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Error converting %s to int: %v\n", name, err)
		return time.Duration(defaultValue)
	}
	return time.Duration(value)
}
func GetDuration(operation string) time.Duration {
	switch operation {
	case "+":
		return time.Duration(TIME_ADDITION_MS)
	case "-":
		return time.Duration(TIME_SUBTRACTION_MS)
	case "*":
		return time.Duration(TIME_MULTIPLICATIONS_MS)
	case "/":
		return time.Duration(TIME_DIVISIONS_MS)
	}
	return time.Duration(1000)
}
