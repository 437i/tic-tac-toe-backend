.PHONY: run build test test-v fmt tidy db-up db-down migrate-up migrate-down ui

run:
	go run ./cmd

build:
	mkdir -p bin
	go build -o bin/tic-tac-toe ./cmd

test:
	go test ./...

test-v:
	go test -v ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

tidy:
	go mod tidy

db-up:
	docker compose up -d

db-down:
	docker compose down

migrate-up:
	goose up

migrate-down:
	goose down

ui:
	python3 manual-test-ui/server.py
