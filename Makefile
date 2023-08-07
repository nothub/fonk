MOD_NAME = $(shell go list -m)
BIN_NAME = $(shell basename $(MOD_NAME))
GOFLAGS  = -race -tags netgo,timetzdata,sqlite_omit_load_extension
LDFLAGS  = -ldflags="-extldflags=-static"

honk: schema.sql $(shell ls go.mod go.sum *.go **/*.go)
	CGO_ENABLED=1 go build $(GOFLAGS) $(LDFLAGS) -o honk

.PHONY: clean
clean:
	go clean
	$(shell docker rm -f "n0thub/honk:latest" 2> /dev/null)

.PHONY: test
test:
	go vet
	go test

.PHONY: image
image:
	docker build -t "n0thub/honk:latest" .
