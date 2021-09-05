FROM golang:alpine AS builder
RUN apk --no-cache add make git

WORKDIR /opt/easemesh-shadowservice-controller
COPY . .

RUN make clean && make build

# ---

FROM alpine:latest

WORKDIR /opt/easemesh-shadowservice-controller

COPY --from=builder /opt/easemesh-shadowservice-controller/bin/ /opt/easemesh-shadowservice-controller/bin/

RUN apk --no-cache add tini tzdata && \
        chmod +x /entrypoint.sh /opt/easemesh-shadowservice-controller/bin/*

ENV PATH /opt/easemesh-shadowservice-controller/bin:$PATH

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
