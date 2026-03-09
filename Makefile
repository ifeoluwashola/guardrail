.PHONY: build test run install

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOINSTALL=$(GOCMD) install

# Binary name
BINARY_NAME=guardrail

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

install:
	$(GOINSTALL)

clean:
	rm -f $(BINARY_NAME)
