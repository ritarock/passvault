PHONY: test build install

BINDIR=bin

test:
	go test ./...

build:
	go build -o ${BINDIR}/passvault .

install:
	go install -o passvault
