PKG       := github.com/ChrisWiegman/kana-cli
VERSION   := $(shell git describe --tags || echo "0.0.1")
GITHASH   := $(shell git describe --tags --always --dirty)
TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

.PHONY: install
install:
	go mod vendor
	go run ./scripts/generateTemplateConstants.go
	go install \
		-ldflags "-s -w -X $(PKG)/internal/cmd.Version=$(VERSION) -X $(PKG)/internal/cmd.GitHash=$(GITHASH) -X $(PKG)/internal/cmd.Timestamp=$(TIMESTAMP)" \
		./cmd/...

.PHONY: update
update:
	go get -u ./...
