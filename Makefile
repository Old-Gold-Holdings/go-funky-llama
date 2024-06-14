.PHONY: start build run

start: build run

build:
	go build -o bin/$(APP_NAME) main.go

run:
	go run main.go