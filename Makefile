PATH := $(PATH):$(HOME)/go/bin
gen:
	protoc -I internal/proto --go_out=paths=source_relative:pkg/api --go-grpc_out=paths=source_relative:pkg/api internal/proto/calculator.proto
run:
	docker-compose up
stop:
	docker-compose stop
update:
	docker-compose down
	docker-compose up --build
runTest:
	go run cmd/orchestrator/main.go
runTest2:
	go run cmd/agent/main.go