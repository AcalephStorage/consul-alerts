FROM golang:1.5
MAINTAINER Acaleph <admin@acale.ph>
RUN apt-get update && apt-get install -y unzip --no-install-recommends && rm -rf /var/lib/apt/lists/*
RUN (curl -OL https://dl.bintray.com/mitchellh/consul/0.5.2_linux_amd64.zip &&\
unzip 0.5.2_linux_amd64.zip &&\
chmod +x consul &&\
mv consul /bin/ &&\
rm 0.5.2_linux_amd64.zip)

ADD . /go/src/github.com/AcalephStorage/consul-alerts
RUN go get github.com/AcalephStorage/consul-alerts

EXPOSE 9000
CMD []
ENTRYPOINT [ "/go/bin/consul-alerts", "--alert-addr=0.0.0.0:9000" ]
