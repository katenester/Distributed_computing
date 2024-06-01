package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/katenester/Distributed_computing/internal/orchestrator/handle"
	"log"
	"net/http"
)

func Run() {
	log.Println("create router")
	// Роутер (маршрутезатор)
	router := httprouter.New()
	log.Println("register user handler")
	handler := handle.NewHandler()
	handler.Register(router)

	log.Println("start application :8080")
	log.Fatalln(http.ListenAndServe(":8080", router))
}
