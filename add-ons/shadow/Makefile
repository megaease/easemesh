SHELL:=/bin/bash

.PHONY: vet clean test dependency_sync image build

MKFILE_PATH := $(abspath $(lastword ${MAKEFILE_LIST}))
MKFILE_DIR := $(shell dirname ${MKFILE_PATH})
ALL_FILES := $(shell find ${MKFILE_DIR}/{cmd,pkg} -type f -name "*.go")

ifndef COMMIT
  COMMIT := $(shell cd ${MKFILE_DIR} && git describe --always HEAD)
endif
ifndef BUILDDATE
  BUILDDATE := $(shell date "+%Y-%m-%d_%H:%M:%SZ+8")
endif

GO_BUILD_FLAGS := -ldflags "-X main.commit=${COMMIT} -X main.buildDate=${BUILDDATE}" -v

TARGET = ${MKFILE_DIR}/bin/easemesh-shadowservice-controller

build: vet ${TARGET}

${TARGET}: ${ALL_FILES}
	@cd ${MKFILE_DIR} && \
	CGO_ENABLED=0 go build ${GO_BUILD_FLAGS}\
  	-o ${TARGET} ${MKFILE_DIR}/cmd/main


vet:
	@cd ${MKFILE_DIR} && go vet ./{cmd,pkg}/...

dependency_sync:
	@cd ${MKFILE_DIR} && go mod tidy -v

clean:
	@rm -f ${TARGET}

image: rootfs/Dockerfile vet ${ALL_FILES}
	docker buildx build --platform linux/amd64 \
	-t megaease/easemesh-shadowservice-controller:latest \
	-f ${MKFILE_DIR}/rootfs/Dockerfile ${MKFILE_DIR} \
	--build-arg COMMITID=${COMMIT} \
	--build-arg BUILDDATE=${BUILDDATE}
	docker tag megaease/easemesh-shadowservice-controller:latest megaease/easemesh-shadowservice-controller:v1.1.0

test:
	go test ./... ${TEST_FLAGS}

