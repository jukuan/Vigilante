.PHONY: build test lint clean run

build:
	go build -o vigilante .

test:
	go test ./... -v

lint:
	golangci-lint run ./...

clean:
	rm -f vigilante state.json vigilante.log vigilante.pid alerts.log

run:
	./vigilante config.yaml

dev-setup:
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	chmod +x bin/*.sh scripts/*.sh
