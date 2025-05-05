PHONY: test build install

BINDIR=bin

test:
	go test ./...

build:
	go build -o ${BINDIR}/pa55 .

install:
	go install -o pa55
