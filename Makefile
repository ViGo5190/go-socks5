all: test vet build
GOFILES=`go list ./... | grep -v vendor`

test:
	go test $(GOFILES)

bench:
	go test -bench=. -benchmem $(GOFILES)

cover:
	go test -coverprofile=cover.out $(GOFILES) && go tool cover -html=cover.out -o cover.html

build: test
	go build .

fmt:
	go fmt $(GOFILES)

vet:
	go vet $(GOFILES)

lint:
	golint $(GOFILES)

docker:
	docker build -t gosocks5 .

travis:  test
	echo "done all"

pre: fmt lint vet test bench