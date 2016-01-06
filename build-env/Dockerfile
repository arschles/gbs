FROM golang:1.5.2

ADD build.sh ./build.sh
RUN mv build.sh /bin && chmod +x /bin/build.sh

ENV GO15VENDOREXPERIMENT=1
ENV CGO_ENABLED=0

CMD /bin/build.sh
