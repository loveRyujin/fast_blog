GO ?= go
EXECUTABLE := fb-apiserver
GOFILES := $(shell find . -type f -name "*.go")

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	$(GO) build -x -o ./_output/$@  ./cmd/$(EXECUTABLE)