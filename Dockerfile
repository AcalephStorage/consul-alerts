FROM alpine:3.3
MAINTAINER psev <paulseverance@gmail.com>

ENV GOPATH /go

RUN mkdir -p /go && \
    apk update && \
    apk add go git ca-certificates && \
    go get -v github.com/AcalephStorage/consul-alerts && \
    mv /go/bin/consul-alerts /bin && \
    rm -rf /go && \
    apk del --purge go git && \
    rm -rf /var/cache/apk/*

EXPOSE 9000
CMD []
ENTRYPOINT [ "/bin/consul-alerts", "--alert-addr=0.0.0.0:9000" ]
