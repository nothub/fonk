honk: schema.sql $(shell ls go.mod go.sum *.go **/*.go)
	go build -race -ldflags="-extldflags=-static" -tags netgo -o honk

.PHONY: clean
clean:
	go clean

.PHONY: test
test:
	go vet
	go test

.PHONY: image
image:
	docker build --no-cache -t "honk:dev" .
