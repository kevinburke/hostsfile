.PHONY: test release

SHELL = /bin/bash -o pipefail

BUMP_VERSION := $(GOPATH)/bin/bump_version
STATICCHECK := $(GOPATH)/bin/staticcheck
RELEASE := $(GOPATH)/bin/github-release
WRITE_MAILMAP := $(GOPATH)/bin/write_mailmap

UNAME := $(shell uname)

test:
	go test ./...

$(STATICCHECK):
	go get honnef.co/go/tools/cmd/staticcheck

$(BUMP_VERSION):
	go get -u github.com/kevinburke/bump_version

$(RELEASE):
	go get -u github.com/aktau/github-release

$(WRITE_MAILMAP):
	go get -u github.com/kevinburke/write_mailmap

force: ;

AUTHORS.txt: force | $(WRITE_MAILMAP)
	$(WRITE_MAILMAP) > AUTHORS.txt

authors: AUTHORS.txt

lint: | $(STATICCHECK)
	$(STATICCHECK) ./...
	go vet ./...

race-test: lint
	go test -race ./...

# Run "GITHUB_TOKEN=my-token make release version=0.x.y" to release a new version.
release: race-test
	$(BUMP_VERSION) minor cmd.go
	git push origin --tags
