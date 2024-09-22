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
        -buildvcs=false \
        -ldflags "-s -w -X github.com/ChrisWiegman/kana-dev/internal/cmd.Version=1.0.0 -X github.com/ChrisWiegman/kana-dev/internal/cmd.Timestamp=2024-03-16_10:50:11PM" \
        ./cmd/... && \
    go test \
        -v \
        -timeout 30s\
        -cover \
        ./...