ready: generate-mock-all lint test-unit test-integration

wire:
	wire ./...

lint:
	golangci-lint run

test-unit:
	go test ./internal/... -race -coverprofile=coverage_unit.out -covermode=atomic -v

test-integration:
	go test -tags integration ./internal/.../handler/... -race -coverprofile=coverage_integration.out -coverpkg=./internal/.../handler/... -covermode=atomic -v

i18n-merge:
	goi18n merge -outdir=./locale ./locale/active.*.toml ./locale/translate.*.toml

generate-mock-all:

