.DEFAULT_GOAL := test

fetch-dependencies:
	go get -v -t -d ./...

test: fetch-dependencies
	go test -v ./...
