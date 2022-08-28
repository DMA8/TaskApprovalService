all: docker.start test.integration docker.stop

docker.start:
	docker-compose -f internal/tests/docker-compose.yaml up -d
	sleep 5

docker.stop:
	docker-compose -f internal/tests/docker-compose.yaml kill

docker.restart: docker.stop docker.start

test.unit:
	go test ./... -cover

test.integration:
	go test -tags=integration ./internal/tests -v -count=1

gen:
	mockgen -source=internal/ports/tasks.go \
	-destination=internal/mocks/mock_tasks.go
	mockgen -source=internal/ports/auth_grpc.go \
	-destination=internal/mocks/mock_auth_grpc.go
	mockgen -source=internal/ports/task_grpc.go \
	-destination=internal/mocks/mock_task_grpc.go
	mockgen -source=internal/ports/tasks_storage.go \
	-destination=internal/mocks/mock_tasks_storage.go

swag:
	swag init -g internal/api/api.go
