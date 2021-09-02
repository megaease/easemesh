FROM golang:alpine AS builder

ARG COMMITID
ARG BUILDDATE

RUN apk add build-base musl-dev
# Install git.
# Git is required for fetching the dependencies.
WORKDIR /easemesh-shadow-service-controller
COPY go.mod go.sum /easemesh-shadow-service-controller/
# Fetch dependencies.
# Using go get.
RUN go mod download

COPY ./pkg/ /easemesh-shadow-service-controller/pkg
COPY ./cmd/ /easemesh-shadow-service-controller/cmd

# Build the binary.
RUN go build  -tags musl \
    -installsuffix cgo -ldflags "-X main.commit=${COMMITID} -X main.buildDate=${BUILDDATE} -extldflags '-static'" \
    -o /easemesh-shadow-service-controller/bin/easemesh-shadow-service-controller \
    github.com/megaease/easemesh/shadow/cmd/main

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /easemesh-shadow-service-controller/bin/easemesh-shadow-service-controller  /easemesh-shadow-service-controller
# Run the hello binary.
CMD ["/easemesh-shadow-service-controller "]
