all: test build

test:
	echo "no test"

build: test
	go build .

fmt:
	go fmt proxy/*.go
	go fmt main.go

pre: fmt test
