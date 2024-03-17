FROM ubuntu:22.04

MAINTAINER Chris Wiegman <contact@chriswiegman.com>

RUN apt-get update && \
    apt-get -qy full-upgrade && \
    apt-get install -qy curl && \
    apt-get install -qy curl && \
    curl -sSL https://get.docker.com/ | sh

RUN curl -OL https://golang.org/dl/go1.22.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xvf go1.22.1.linux-amd64.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"

CMD go build \
        -o ./build/kana \
        -ldflags "-s -w -X $(PKG)/internal/cmd.Version=$(VERSION) -X $(PKG)/internal/cmd.Timestamp=$(TIMESTAMP)" \
        ./cmd/... && \
    go test \
        -v \
        -timeout 30s\
        -cover \
        ./...