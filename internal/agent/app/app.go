package app

import (
	"github.com/katenester/Distributed_computing/internal/agent/services"
	"log"
	"os"
	"strconv"
	"sync"
)

func Run() {
	log.Println("Запуск агента")
	//Количество горутин регулируется переменной среды	COMPUTING_POWER
	// Получаем значение переменной среды COMPUTING_POWER
	cpuCount, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || cpuCount <= 0 {
		cpuCount = 3 // Используем значение по умолчанию, если переменная не задана или содержит некорректное значение
	}
	wg := sync.WaitGroup{}
	wg.Add(cpuCount)
	for i := 0; i < cpuCount; i++ {
		// Запускаем вычислительные машины в горутинах
		go services.Demon()
	}
	wg.Wait()
}
