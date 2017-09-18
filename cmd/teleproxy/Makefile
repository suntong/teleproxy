SHELL       = /bin/bash
STAMP      ?= $$(date +%Y-%m-%d_%H:%M.%S)
APP        ?= $(shell basename $$PWD)
APP_DEPS    = *.go
BUILD_FLAG  =

# Search .git for commit id fetch
GIT_ROOT  ?= $$([ -d ./.git ] && echo "." || { [ -d ../.git ] && echo ".." ; } || { [ -d ../../.git ] && echo "../.." ; })

all: build

build: vet $(APP)
#build: lint vet $(APP)

# build app
$(APP): vendor $(APP_DEPS)
	@echo "*** $@ ***"
	[ -d $(GIT_ROOT)/.git ] && GH=`git rev-parse HEAD` || GH=nogit ; \
	go build $$BUILD_FLAG -v -ldflags "-X main.Build=$(STAMP) -X main.Commit=$$GH" -o $(APP) $(APP_DEPS)

# install vendor deps
vendor:
	@echo "*** $@ ***"
	which glide > /dev/null || curl https://glide.sh/get | sh
	@echo "*** $@:glide ***"
	glide install

# clean binary
clean:
	@[ -f $(APP) ] && rm $(APP) || true
	@[ -d vendor ] && rm -rf vendor || true

## run go lint
lint:
	@echo "*** $@ ***"
	@golint *.go

## run go vet
vet:
	@echo "*** $@ ***"
	@go vet *.go
