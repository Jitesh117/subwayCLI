APP_NAME := subwaycli
MAIN := ./cmd/subwaycli

.PHONY: run build build-all clean fmt test

run:
	go run $(MAIN)

build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) $(MAIN)

build-all:
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build -o dist/$(APP_NAME)-darwin-amd64 $(MAIN)
	GOOS=darwin GOARCH=arm64 go build -o dist/$(APP_NAME)-darwin-arm64 $(MAIN)
	GOOS=linux GOARCH=amd64 go build -o dist/$(APP_NAME)-linux-amd64 $(MAIN)
	GOOS=linux GOARCH=arm64 go build -o dist/$(APP_NAME)-linux-arm64 $(MAIN)
	GOOS=windows GOARCH=amd64 go build -o dist/$(APP_NAME)-windows-amd64.exe $(MAIN)
	GOOS=windows GOARCH=arm64 go build -o dist/$(APP_NAME)-windows-arm64.exe $(MAIN)

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -rf bin dist
