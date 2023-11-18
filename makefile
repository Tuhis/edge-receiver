.PHONY: build run dev

build:
	go build -o edge-receiver cmd/edge-receiver/main.go

run:
	./edge-receiver

dev:
	reflex -r '\.go$$' -s -- sh -c 'go run cmd/edge-receiver/main.go'

test:
	go test -v ./...

test-cover:
	go test -cover ./...
