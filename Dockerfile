FROM golang:1.12-alpine
MAINTAINER tselikov@uchi.ru
ENV GOPATH /go

RUN mkdir -p /go && \
    apk update && \
    apk add bash ca-certificates git && \
    GO111MODULE="off" go get -v github.com/uchiru/consul-alerts && \
    mv /go/bin/consul-alerts /bin && \
    GO111MODULE="off" go get -v github.com/hashicorp/consul && \
    mv /go/bin/consul /bin && \
    rm -rf /go && \
    apk del --purge go git alpine-sdk && \
    rm -rf /var/cache/apk/*]

FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=0 /bin /bin

EXPOSE 9000
CMD []
ENTRYPOINT [ "/bin/consul-alerts", "--alert-addr=0.0.0.0:9000" ]
