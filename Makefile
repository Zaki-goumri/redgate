APP=redgate

.PHONY: build run test tidy docker-up docker-down

build:
	go build -o bin/$(APP) ./cmd/redgate

run:
	go run ./cmd/redgate --config ./config.example.yaml

test:
	go test ./...

tidy:
	go mod tidy

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v
