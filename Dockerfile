FROM alpine:3.1

RUN apk add -U docker && rm -rf /var/cache/apk/*

ADD build.sh ./build.sh
ADD gbs ./gbs
# note that currently gbs has to be in the same dir as build.sh
RUN mv build.sh /bin/build.sh && mv gbs /bin/gbs
CMD gbs
