BINDIR := $(CURDIR)/bin
BINNAME ?= kubesweeper

# go option
PKG := ./...
TAGS :=
TESTS := .
TESTFLAGS :=
LDFLAGS := -w -s
GOFLAGS :=
SRC  := $(shell find . -type f -name '*.go' -print)

.PHONY: all
all: build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd

.PHONY: test
test:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)
