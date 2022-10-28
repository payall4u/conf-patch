
GOOS ?= linux
GOARCH ?= amd64

.PHONY: all test clean

all: clean test binary

GOFILES = $(shell find . -name '*.go')

test:
	go test .

binary:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o _output/$(GOOS)/$(GOARCH)/conf-patch -a -ldflags '-extldflags "-static"' .

clean:
	rm -rf _output/*