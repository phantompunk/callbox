.PHONY: all build run test clean watch

all: build test

build:
	@go build -o tmp/main cmd/api/main.go

run:
	@go run cmd/api/main.go

test:
	@go test -v ./...

clean:
	@rm -f tmp/main

watch:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Go's 'air' is not installed." \
		echo "Install using 'go install github.com/air-verse/air@latest"; \
	fi
