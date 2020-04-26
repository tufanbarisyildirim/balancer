PWD := $(shell pwd)
export GO111MODULE=on

test:
	go test -race -cover ${PWD}

bench:
	go test -bench=.
	
fmt:
	find . -name "*.go" | xargs gofmt -w -s

deps:
	go get -v all

.PHONY: fmt deps test