FROM alpine:3.13

WORKDIR /opt/easemesh-shadowservice

ADD bin/easemesh-shadowservice-controller /opt/easemesh-shadowservice/bin/
COPY rootfs/entrypoint.sh /

RUN apk --no-cache add tini tzdata && \
        chmod +x /entrypoint.sh /opt/easemesh-shadowservice/bin/*

ENV PATH /opt/easemesh-shadowservice/bin:$PATH

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
