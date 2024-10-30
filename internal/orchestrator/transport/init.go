package transport

import "github.com/katenester/Distributed_computing/internal/orchestrator/Storage"

var storage *Storage.Storage

// Инициализируем глобальную переменную
func init() {
	storage = Storage.New()
}
