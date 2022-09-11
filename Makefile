PKG       := github.com/ChrisWiegman/kana-cli/internal/cmd.
VERSION   := $(shell git describe --tags || echo "0.0.1")
GITHASH   := $(shell git describe --tags --always --dirty)
TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

.PHONY: install
install:
	go mod vendor
	go run ./scripts/generateTemplateConstants.go
	go install \
		-ldflags "-s -w -X $(PKG)Version=$(VERSION) -X $(PKG)GitHash=$(GITHASH) -X $(PKG)Timestamp=$(TIMESTAMP)" \
		./cmd/...

.PHONY: update
update:
	go get -u ./...
