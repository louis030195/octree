GOFILES = $(shell find . -name '*.go')

default: build

build:
	mkdir -p build

build: build/octree

build-native: $(GOFILES)
	go build -o build/native-octree .

build/octree: $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/octree .