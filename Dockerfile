FROM google/golang:1.4
MAINTAINER Acaleph <admin@acale.ph>

WORKDIR /gopath/src/github.com/AcalephStorage/consul-alerts
ADD . /gopath/src/github.com/AcalephStorage/consul-alerts
RUN go get consul-alerts

EXPOSE 9000
CMD []
ENTRYPOINT [ "/gopath/bin/consul-alerts", "--alert-addr=0.0.0.0:9000" ]
