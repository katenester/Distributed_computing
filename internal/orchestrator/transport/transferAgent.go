package transport

import (
	"context"
	"errors"
	proto "github.com/katenester/Distributed_computing/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type GRPSServer struct {
	proto.UnimplementedGenerateTaskServer
}

// GiveTask - Передача задачи агенту
func (g *GRPSServer) GiveTask(context.Context, *emptypb.Empty) (*proto.TaskResponse, error) {
	// Ищем задачу для агента
	log.Println("Зашли за таской")
	if task, isTask := storage.FindTask(); isTask {
		log.Println("тАСКА", &proto.TaskResponse{
			Id:         int32(task.Id),
			X:          float32(task.Arg1),
			Y:          float32(task.Arg2),
			Operator:   task.Operation,
			LastAccess: timestamppb.New(task.LastAccess),
		})
		return &proto.TaskResponse{
			Id:         int32(task.Id),
			X:          float32(task.Arg1),
			Y:          float32(task.Arg2),
			Operator:   task.Operation,
			LastAccess: timestamppb.New(task.LastAccess)}, nil
	}
	// Если не нашлась задача => возвращаем gRPC код NotFound
	return nil, status.Error(codes.NotFound, "task not found")
}

// GetResult - Получение результата вычисления от агента
func (g *GRPSServer) GetResult(ctx context.Context, req *proto.ResultRequest) (*emptypb.Empty, error) {
	log.Println("Получили результат таски", req)
	switch res := req.Details.(type) {
	case *proto.ResultRequest_Result:
		//Проверяем таску в хранилке , и если такая есть => меняем результат
		storage.FindAndReplace(req.GetId(), float64(req.GetResult()), nil)
	case *proto.ResultRequest_Error:
		if res.Error == "division by zero" {
			// поиск задачи и сделать выражение ошибочным
			storage.FindAndReplace(req.GetId(), 0.0, errors.New(req.GetError()))
		}
	default:
		log.Println("Unknown result type")
		// Обработка случая, когда ни одно поле не установлено
	}
	return nil, nil
}

// GiveTask - Передача задачи агенту
//func (h *handler) GiveTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//	// Ищем задачу для агента
//	if task, isTask := storage.FindTask(); isTask {
//		resp := map[string]model.Task{
//			"task": task,
//		}
//		log.Println("Сервер: Передаю задачу агенту", resp)
//		w.Header().Set("Content-Type", "application/json")
//		err := json.NewEncoder(w).Encode(resp)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusCreated)
//		return
//	}
//	// Задачи нет
//	w.WriteHeader(http.StatusNotFound)
//}

// GetResultTask - Получение результата вычисления от агента
//func (h *handler) GetResultTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//	// Получаем id и result
//	data := struct {
//		Id     int         `json:"id"`
//		Result float64     `json:"result"`
//		Err    interface{} `json:"error"`
//	}{}
//	err := json.NewDecoder(r.Body).Decode(&data)
//	if err != nil {
//		w.WriteHeader(http.StatusUnprocessableEntity)
//		return
//	}
//	log.Println("Сервер: Получение результата вычисления от агента", data)
//	// Проверяем, не была ли ошибка в вычислениях(деление на 0)
//	if data.Err == nil {
//		err = nil
//	} else {
//		err = errors.New(fmt.Sprint(data.Err))
//	}
//	// Проверяем таску в хранилке , и если такая есть => меняем результат
//	if storage.FindAndReplace(data.Id, data.Result, err) {
//		w.WriteHeader(http.StatusCreated)
//		return
//	}
//	w.WriteHeader(http.StatusNotFound)
//}
