.PHONY: build test vet check

build:
	go build -o bin/zigguard ./cmd/zigguard

test:
	go test ./...

vet:
	go vet ./...

check: test vet
