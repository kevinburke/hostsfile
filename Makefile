.PHONY: test

vet: 
	go vet ./...

test: vet
	go test ./...

race-test: vet
	go test -race ./...
