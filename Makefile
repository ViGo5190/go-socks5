all: test build

test:
	go test ./...

cover:
	go test -coverprofile=cover.out ./... && go tool cover -html=cover.out -o cover.html

build: test
	go build .

fmt:
	go fmt proxy/*.go
	go fmt main.go

pre: fmt test
