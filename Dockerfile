FROM golang:alpine AS build-env
MAINTAINER Jason Berlinsky <jason@barefootcoders.com>

RUN echo "@edge http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
  apk add --update --no-cache \
    make \
    git \
    wget \
    ca-certificates \
    openssl \
    glide@edge

RUN mkdir -p /go/src/github.com/delectable/gosubscriber_redux
WORKDIR /go/src/github.com/delectable/gosubscriber_redux

ADD . /go/src/github.com/delectable/gosubscriber_redux

RUN make installdeps
RUN make clean
RUN make
