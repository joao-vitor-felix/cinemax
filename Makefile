build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go -migrate

fmt:
	golangci-lint fmt

lint:
	golangci-lint run

tidy:
	go mod tidy

.PHONY: test
test:
	go test -v ./...

.PHONY: docs
docs:
	swag init -g cmd/api/main.go
