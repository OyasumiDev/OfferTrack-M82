.PHONY: build run dev test lint clean release

build:
	go build -o ../bin/offertrack ./cmd/offertrack

dev:
	go run -race ./cmd/offertrack

run: build
	../bin/offertrack

test:
	go test -v -race ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	go vet ./...

clean:
	rm -f ../bin/offertrack coverage.out

release:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/offertrack ./cmd/offertrack
