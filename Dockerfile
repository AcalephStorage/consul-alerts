FROM progrium/busybox:latest
MAINTAINER Acaleph <admin@acale.ph>

ADD https://dl.bintray.com/mitchellh/consul/0.4.1_linux_amd64.zip /tmp/consul.zip
RUN cd /bin && unzip /tmp/consul.zip && chmod +x /bin/consul && rm /tmp/consul.zip

ADD http://dl.bintray.com/darkcrux/generic/consul-alerts-0.1.1-linux-amd64.tar /tmp/consul-alerts.tar
RUN cd /bin && tar -xf /tmp/consul-alerts.tar && chmod +x /bin/consul-alerts && rm /tmp/consul-alerts.tar

EXPOSE 9000
ENTRYPOINT [ "bin/consul-alerts", "--alert-addr=0.0.0.0:9000" ]
CMD []
