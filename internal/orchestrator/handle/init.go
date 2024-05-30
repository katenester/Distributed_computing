package handle

import "main.go/internal/orchestrator/Storage"

var storage *Storage.Storage

// Инициализируем глобальную переменную
func init() {
	storage = Storage.New()
}
