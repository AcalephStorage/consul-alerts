FROM alpine:edge
MAINTAINER tselikov@uchi.ru
ENV GOPATH /go

RUN mkdir -p /go && \
    apk update && \
    apk add bash ca-certificates git go alpine-sdk && \
    go get -v github.com/uchiru/consul-alerts && \
    mv /go/bin/consul-alerts /bin && \
    go get -v github.com/hashicorp/consul && \
    mv /go/bin/consul /bin && \
    rm -rf /go && \
    apk del --purge go git alpine-sdk && \
    rm -rf /var/cache/apk/*

EXPOSE 9000
CMD []
ENTRYPOINT [ "/bin/consul-alerts", "--alert-addr=0.0.0.0:9000" ]
