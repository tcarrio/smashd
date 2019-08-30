PKGS := $(shell go list ./...)
.PHONY: test
test: lint
	go test $(PKGS)

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor

BINARY := powerd
VERSION := $(shell go run cmd/main.go -v)
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64 cmd/main.go

.PHONY: raspberrypi
raspberrypi:
	mkdir -p release
	GOOS=linux GOARCH=arm go build -o release/$(BINARY)-$(VERSION)-linux-arm cmd/main.go

.PHONY: release
release: windows linux darwin raspberrypi
