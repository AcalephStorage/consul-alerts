FROM google/golang:1.4
MAINTAINER Acaleph <admin@acale.ph>

WORKDIR /gopath/src/consul-alerts
ADD . /gopath/src/consul-alerts
RUN go get consul-alerts

EXPOSE 9000
CMD []
ENTRYPOINT [ "/gopath/bin/consul-alerts", "--alert-addr=0.0.0.0:9000" ]
