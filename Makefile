# Go parameters
GOCMD=go
BUILD_DIR=./build
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
PORT=8181

all: test build

start: build run

build: clean
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

run:
	$(BUILD_DIR)/$(BINARY_NAME) build --input pages/ --output build/web/
	$(BUILD_DIR)/$(BINARY_NAME) serve --input build/web/ --port $(PORT)

deps:
	$(GOGET) -v ./...
