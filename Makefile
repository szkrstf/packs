GIT_COMMIT = $(shell git rev-parse --short HEAD || echo 'dev')
VERSION = $(shell git describe --tags || echo 'dev')
GOLDFLAGS="-X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT)"
APP=packs

all: build

clean:
	rm -f $(APP)

test:
	go test -v ./...

build: clean test
	go build -ldflags=$(GOLDFLAGS) -o $(APP) ./cmd/$(APP)

build-linux: clean test
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags=$(GOLDFLAGS) -o $(APP) ./cmd/$(APP)

.PHONY: all clean test build build-linux