.PHONY: run build tidy test swag docker

run:
	go run ./cmd/pureblog

build:
	CGO_ENABLED=0 go build -o bin/pureblog ./cmd/pureblog

tidy:
	go mod tidy

test:
	go test ./...

# Generate OpenAPI (docs/swagger.json) consumed by frontend api-types codegen.
swag:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/pureblog/main.go -o docs --parseInternal

docker:
	docker build -t pureblog-backend .
