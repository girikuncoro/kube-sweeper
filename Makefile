BINDIR := $(CURDIR)/bin
NAME ?= kubesweeper
VERSION ?= 0.1.0

# go option
GO ?= go
PKG := ./...
TAGS :=
TESTS := .
TESTFLAGS :=
LDFLAGS := -w -s
GOFLAGS :=
SRC  := $(shell find . -type f -name '*.go' -print)

# docker option
REGISTRY ?= girikuncoro

.PHONY: all
all: build

.PHONY: build
build: $(BINDIR)/$(NAME)

$(BINDIR)/$(NAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(NAME) ./cmd

.PHONY: test
test:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

.PHONY: docker
docker:
	@echo
	@echo "==> Building docker image <=="
	docker build -t $(REGISTRY)/$(NAME):$(VERSION) .

.PHONY: docker-publish
docker-publish:
	@echo
	@echo "==> Publishing docker image <=="
	docker push $(REGISTRY)/$(NAME):$(VERSION)
