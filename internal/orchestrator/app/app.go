package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/katenester/Distributed_computing/internal/orchestrator/transport"
	proto "github.com/katenester/Distributed_computing/pkg/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func Run() {
	log.Println("create router")
	// Роутер (маршрутезатор)
	router := httprouter.New()
	log.Println("register user handler")
	handler := transport.NewHandler()
	handler.Register(router)
	// Запускаем gRPC сервер в отдельной горутине
	go startGRPCServer()

	log.Println("start application :" + viper.GetString("http_port"))
	log.Fatalln(http.ListenAndServe(viper.GetString("http_port"), router))
}

func startGRPCServer() {
	listen, err := net.Listen("tcp", viper.GetString("grpc_port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// Регистрируем сервисы gRPC
	proto.RegisterGenerateTaskServer(s, &transport.GRPSServer{})

	log.Println("gRPC server is running on :", viper.GetString("grpc_port"))
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
