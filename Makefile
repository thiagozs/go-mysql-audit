BINARY=paudit


VERSION=$(shell git describe --tags)
BUILD=$(shell git rev-parse --short HEAD)

ifeq (${VERSION},)
VERSION="beta"
endif

LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

build:
	go build ${LDFLAGS} -o ${BINARY} main.go
