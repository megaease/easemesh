
SHELL:=/bin/bash
.PHONY: build fmt vet clean \
		mod_update vendor_from_mod vendor_clean test

# Path Related
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))

# Version
RELEASE?=1.0.0
# Git Related
GIT_REPO_INFO=$(shell cd ${MKFILE_DIR} && git config --get remote.origin.url)
ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif

ALL_FILES = $(shell find ${MKFILE_DIR}{cmd,pkg} -type f -name "*.go")

GO_LD_FLAGS= "-s -w -X github.com/megaease/easemeshctl/pkg/version.RELEASE=${RELEASE} -X github.com/megaease/easemeshctl/pkg/version.COMMIT=${GIT_COMMIT} -X github.com/megaease/easemeshctl/pkg/version.REPO=${GIT_REPO_INFO}"

TARGET=${MKFILE_DIR}bin/emctl

all: build

test:
	@go list ${MKFILE_DIR}/cmd/... | grep -v -E 'vendor' | xargs -n1 go test ${TEST_FLAGS}

clean:
	rm -rf ${TARGET}

fmt:
	cd ${MKFILE_DIR} && go fmt ./cmd/...

vet:
	cd ${MKFILE_DIR} && go vet ./cmd/...

vendor_from_mod:
	cd ${MKFILE_DIR} && go mod vendor

vendor_clean:
	rm -rf ${MKFILE_DIR}vendor

mod_update:
	cd ${MKFILE_DIR} && go get -u

${TARGET}: ${ALL_FILES}
	@echo "build emctl "
	cd ${MKFILE_DIR} && \
	CGO_ENABLED=0 go build -v -ldflags ${GO_LD_FLAGS} \
	-o ${TARGET} ${MKFILE_DIR}cmd/client/main.go

build: ${TARGET}
