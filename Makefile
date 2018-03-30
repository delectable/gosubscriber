GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
RELEASE_DIR=bin
DEVTOOL_DIR=devtools
LIBRARY_NAME=gosubscriber_redux
PACKAGE=github.com/delectable/$(LIBRARY_NAME)

.PHONY: clean build build-linux-amd64 build-linux-386 build-darwin-amd64 build-darwin-386 build-windows-amd64 build-windows-386 $(RELEASE_DIR)/_$(GOOS)_$(GOARCH) all

all: installdeps clean build-linux-amd64 # build-linux-386 build-darwin-amd64 build-darwin-386 build-windows-amd64 build-windows-386

build: $(RELEASE_DIR)/$(PACKAGE)_$(GOOS)_$(GOARCH)

build-linux-amd64:
	$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-386:
	@$(MAKE) build GOOS=linux GOARCH=386

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-darwin-386:
	@$(MAKE) build GOOS=darwin GOARCH=386

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64

build-windows-386:
	@$(MAKE) build GOOS=windows GOARCH=386

$(RELEASE_DIR)/$(PACKAGE)_$(GOOS)_$(GOARCH):
	go build \
		-o $(RELEASE_DIR)/$(PACKAGE)_$(GOOS)_$(GOARCH) gosubscriber.go

installdeps:
	@PATH=$(DEVTOOL_DIR)/$(GOOS)/$(GOARCH):$(PATH) glide install

test:
	go test -v ./... -queues="test_queue"

clean:
	rm -rf $(RELEASE_DIR)/$(PACKAGE)_*
