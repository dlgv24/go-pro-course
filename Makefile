WEEK ?= 01

.PHONY: run lint

run:
	go run ./week$(WEEK)

lint:
	golangci-lint run ./...