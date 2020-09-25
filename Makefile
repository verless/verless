TAG := $(shell git describe --tags --abbrev=0)
COMMIT := $(shell git rev-parse --short HEAD)
PROJECTNAME := verless
TARGET := target
GOFILES := ./cmd/verless

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X github.com/verless/verless/config.GitTag=$(TAG) -X github.com/verless/verless/config.GitCommit=$(COMMIT)"

## build: Compile the binary.
build:
	@mkdir -p $(TARGET)
	@GOPATH=$(GOPATH) \
	GOBIN=$(GOBIN) \
	GOARM=$(GOARM) \
	CGO_CPPFLAGS=$(CGO_CPPFLAGS) \
	CGO_CFLAGS=$(CGO_CFLAGS) \
	CGO_CXXFLAGS=$(CGO_CXXFLAGS) \
	CGO_LDFLAGS=$(CGO_LDFLAGS) \
	go build $(LDFLAGS) -o $(TARGET) $(GOFILES)

## clean the build folder
clean:
	@rm -Rf .target

test:
	@go test ./...

all: build

.PHONY: all build test clean