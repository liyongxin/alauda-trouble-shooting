unexport GOBIN

GO           ?= go
GOFMT        ?= $(GO)fmt
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
GOVENDOR     := $(FIRST_GOPATH)/bin/govendor
pkgs          = ./...

PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
DOCKER_IMAGE_NAME       ?= alauda-trouble-shooting

.PHONY: all
all: style unused build test

.PHONY: style
style:
	@echo ">> checking code style"
	! $(GOFMT) -d $$(find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

.PHONY: test
test:
	@echo ">> running all tests"
	$(GO) test -race $(pkgs)

.PHONY: format
format:
	@echo ">> formatting code"
	$(GO) fmt $(pkgs)

.PHONY: unused
unused: $(GOVENDOR)
	@echo ">> running check for unused packages"
	@$(GOVENDOR) list +unused | grep . && exit 1 || echo 'No unused packages'

.PHONY: docker
docker: build
	docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

.PHONY: build
build:
	@echo ">> building binaries"
	$(GO) build -v -o alauda_trouble_shooting

.PHONY: $(GOVENDOR)
$(GOVENDOR):
	GOOS= GOARCH= $(GO) get -u github.com/kardianos/govendor
