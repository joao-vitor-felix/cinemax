build:
	go build -o bin/api cmd/api/main.go

build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/api cmd/api/main.go

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
