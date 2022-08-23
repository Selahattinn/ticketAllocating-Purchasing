ready: generate-doc generate-mock-all lint test-unit test-integration

wire:
	wire ./...

lint:
	golangci-lint run

test-unit:
	go test ./internal/... -race -coverprofile=coverage_unit.out -covermode=atomic -v

test-integration:
	go test -tags integration ./internal/.../handler/... -race -coverprofile=coverage_integration.out -coverpkg=./internal/.../handler/... -covermode=atomic -v

generate-mock-all:
	mockgen -source=./pkg/mysql/mysql.go -destination=./pkg/mysql/mysql/mysql_mock.go -package=mocks
	mockgen -source=./internal/api/ticket/repository.go -destination=./internal/api/ticket/mocks/repository_mock.go -package=mocks
	mockgen -source=./internal/api/ticket/service.go -destination=./internal/api/ticket/mocks/service_mock.go -package=mocks

generate-doc:
	swag init -g cmd/api/server.go -o docs