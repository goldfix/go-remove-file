# .PHONY: runtime

VERSION = $(shell git describe --tags --abbrev=0)
HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell go run tools/build-date.go)

# Builds grm after checking dependencies but without updating the runtime
build: deps lib
	go build -ldflags "-X main.Version=$(VERSION) -X main.CommitHash=$(HASH) -X 'main.CompileDate=$(DATE)'" ./cmd/grm

# Builds grm after building the runtime and checking dependencies
build-all: runtime build

# Builds grm without checking for dependencies
build-quick:
	go build -ldflags "-X main.Version=$(VERSION) -X main.CommitHash=$(HASH) -X 'main.CompileDate=$(DATE)'" ./cmd/grm

# Same as 'build' but installs to $GOPATH/bin afterward
install: build
	mv grm $(GOPATH)/bin

# Same as 'build-all' but installs to $GOPATH/bin afterward
install-all: runtime install

# Same as 'build-quick' but installs to $GOPATH/bin afterward
install-quick: build-quick
	mv grm $(GOPATH)/bin

# Updates lib
lib:
	git -c $(GOPATH)/src/github.com/satori/go.uuid pull

# Checks for dependencies
deps:
	go get -d ./cmd/grm

# Builds the runtime
runtime:
	go get -u github.com/satori/go.uuid/...
	# $(GOPATH)/bin/go-bindata -nometadata -o runtime.go runtime/...
	# mv runtime.go grm

test:
	rm -rf /tmp/folder_test/
	unzip ./folder_test/folder_test.zip -d /tmp/
	go get -d ./cmd/grm
	go test ./cmd/grm

clean:
	rm -f grm