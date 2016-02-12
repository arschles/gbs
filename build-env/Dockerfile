FROM quay.io/deis/go-dev:0.5.0

ADD build.sh ./build.sh
RUN mv build.sh /bin && \
    chmod +x /bin/build.sh && \
    go get -u github.com/mitchellh/gox

ENV GO15VENDOREXPERIMENT=1
ENV CGO_ENABLED=0

CMD /bin/build.sh
