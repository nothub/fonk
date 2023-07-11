honk: .preflightcheck schema.sql $(shell ls go.mod go.sum *.go **/*.go)
	go build -race -o honk

.preflightcheck: tools/preflight.sh
	@sh ./tools/preflight.sh

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
