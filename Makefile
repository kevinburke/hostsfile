.PHONY: test

SHELL = /bin/bash -o pipefail

BAZEL_VERSION := 0.7.0
BAZEL_DEB := bazel_$(BAZEL_VERSION)_amd64.deb

BUMP_VERSION := $(GOPATH)/bin/bump_version
MEGACHECK := $(GOPATH)/bin/megacheck
RELEASE := $(GOPATH)/bin/github-release
WRITE_MAILMAP := $(GOPATH)/bin/write_mailmap

UNAME := $(shell uname)

test: lint
	bazel test --test_output=errors //...

$(MEGACHECK):
ifeq ($(UNAME), Darwin)
	curl --silent --location --output $(MEGACHECK) https://github.com/kevinburke/go-tools/releases/download/2017-10-04/megacheck-darwin-amd64
endif
ifeq ($(UNAME), Linux)
	curl --silent --location --output $(MEGACHECK) https://github.com/kevinburke/go-tools/releases/download/2017-10-04/megacheck-linux-amd64
endif
	chmod 755 $(MEGACHECK)

$(BUMP_VERSION):
	go get -u github.com/Shyp/bump_version

$(RELEASE):
	go get -u github.com/aktau/github-release

$(WRITE_MAILMAP):
	go get -u github.com/kevinburke/write_mailmap

force: ;

AUTHORS.txt: force | $(WRITE_MAILMAP)
	$(WRITE_MAILMAP) > AUTHORS.txt

authors: AUTHORS.txt

lint: | $(MEGACHECK)
	$(MEGACHECK) ./...
	go vet ./...

race-test: lint
	bazel test --features=race --test_output=errors //...

install-travis:
	wget "https://storage.googleapis.com/bazel-apt/pool/jdk1.8/b/bazel/$(BAZEL_DEB)"
	sudo dpkg --force-all -i $(BAZEL_DEB)
	sudo apt-get install moreutils -y

ci:
	bazel --batch --host_jvm_args=-Dbazel.DigestFunction=SHA1 test \
		--experimental_repository_cache="$$HOME/.bzrepos" \
		--spawn_strategy=remote \
		--test_output=errors \
		--strategy=Javac=remote \
		--noshow_progress \
		--noshow_loading_progress \
		--features=race //... 2>&1 | ts '[%Y-%m-%d %H:%M:%.S]'

# Run "GITHUB_TOKEN=my-token make release version=0.x.y" to release a new version.
release: race-test | $(BUMP_VERSION) $(RELEASE)
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	$(BUMP_VERSION) --version=$(version) cmd.go
	git push origin --tags
	mkdir -p releases/$(version)
	# Change the binary names below to match your tool name
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/hostsfile-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/hostsfile-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/hostsfile-windows-amd64 .
	# Change the Github username to match your username.
	# These commands are not idempotent, so ignore failures if an upload repeats
	$(RELEASE) release --user kevinburke --repo hostsfile --tag $(version) || true
	$(RELEASE) upload --user kevinburke --repo hostsfile --tag $(version) --name hostsfile-linux-amd64 --file releases/$(version)/hostsfile-linux-amd64 || true
	$(RELEASE) upload --user kevinburke --repo hostsfile --tag $(version) --name hostsfile-darwin-amd64 --file releases/$(version)/hostsfile-darwin-amd64 || true
	$(RELEASE) upload --user kevinburke --repo hostsfile --tag $(version) --name hostsfile-windows-amd64 --file releases/$(version)/hostsfile-windows-amd64 || true
