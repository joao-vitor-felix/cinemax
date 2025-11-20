build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

vet:
	go vet ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run
	staticcheck ./...

code-quality:
	make vet
	make fmt
	make lint

tidy:
	go mod tidy

.PHONY: test
test:
	go test -v ./...

.PHONY: docs
docs:
	swag init -g cmd/api/main.go -o ./docs

