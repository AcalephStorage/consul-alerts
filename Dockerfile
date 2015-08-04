FROM google/golang:1.4
MAINTAINER Acaleph <admin@acale.ph>
RUN apt-get install unzip
ADD https://dl.bintray.com/mitchellh/consul/0.5.2_linux_amd64.zip /tmp/consul.zip
RUN cd /bin && unzip /tmp/consul.zip && chmod +x /bin/consul && rm /tmp/consul.zip

WORKDIR /gopath/src/consul-alerts
ADD . /gopath/src/consul-alerts
RUN go get consul-alerts

EXPOSE 9000
CMD []
ENTRYPOINT [ "/gopath/bin/consul-alerts", "--alert-addr=0.0.0.0:9000" ]
