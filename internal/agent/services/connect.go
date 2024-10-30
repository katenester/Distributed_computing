package services

import (
	"context"
	"github.com/katenester/Distributed_computing/internal/model"
	proto "github.com/katenester/Distributed_computing/pkg/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

// Demon - тянет за ручку сервер
func Demon() {
	for {
		log.Println("Дергает агент функцию")
		resp, err := MakeRequestGiveTask()
		log.Println("Получаем", resp)
		// Если задача для вычисления не найдена => скип
		if status.Code(err) == codes.NotFound {
			continue
		}
		if err != nil {
			log.Printf("error grpc client GiveTask %s", err.Error())
			continue
		}
		// Если задача найдена
		task := model.Task{
			Id:         resp.Id,
			Arg1:       float64(resp.X),
			Arg2:       float64(resp.Y),
			Operation:  resp.Operator,
			LastAccess: time.Unix(resp.LastAccess.Seconds, int64(resp.LastAccess.Nanos)),
		}
		log.Println("Конвертируем resopnse", task)
		result, err := calculation(task)
		_, err = MakeRequestGetResult(task.Id, result, err)
		if err != nil {
			log.Printf("error grpc client GetResult %s", err.Error())
			continue
		}
	}
}

func MakeRequestGiveTask() (*proto.TaskResponse, error) {
	conn, err := grpc.Dial("orh"+viper.GetString("grpc_port"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	client := proto.NewGenerateTaskClient(conn)
	return client.GiveTask(ctx, &emptypb.Empty{})
}

func MakeRequestGetResult(idTask int32, result float64, errResult error) (*emptypb.Empty, error) {
	conn, err := grpc.Dial("orh"+viper.GetString("grpc_port"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	client := proto.NewGenerateTaskClient(conn)
	var resultProto *proto.ResultRequest
	if errResult != nil {
		resultProto = &proto.ResultRequest{
			Id: int32(idTask),
			Details: &proto.ResultRequest_Error{ // Устанавливаем поле error
				Error: errResult.Error(),
			},
		}
	} else {
		resultProto = &proto.ResultRequest{
			Id: int32(idTask),
			Details: &proto.ResultRequest_Result{ // Устанавливаем поле error
				Result: float32(result),
			},
		}
	}
	return client.GetResult(ctx, resultProto)
}
