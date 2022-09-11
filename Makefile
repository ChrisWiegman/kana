PKG       := github.com/ChrisWiegman/kana-cli
VERSION   := $(shell git describe --tags || echo "0.0.1")
GITHASH   := $(shell git describe --tags --always --dirty)
TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

.PHONY: change
change:
	docker run \
		--mount type=bind,source=$(PWD),target=/src \
		-w /src \
		-it \
		ghcr.io/miniscruff/changie \
		new

.PHONY: changelog
changelog:
	docker run \
		--mount type=bind,source=$(PWD),target=/src \
		-w /src \
		-it \
		ghcr.io/miniscruff/changie \
		batch $(VERSION)
	docker run \
		--mount type=bind,source=$(PWD),target=/src \
		-w /src \
		-it \
		ghcr.io/miniscruff/changie \
		merge

.PHONY: clean
clean:
	rm -rf \
		dist \
		vendor

.PHONY: install
install:
	go mod vendor
	go run ./scripts/generateTemplateConstants.go
	go install \
		-ldflags "-s -w -X $(PKG)/internal/cmd.Version=$(VERSION) -X $(PKG)/internal/cmd.GitHash=$(GITHASH) -X $(PKG)/internal/cmd.Timestamp=$(TIMESTAMP)" \
		./cmd/...

.PHONY: release
release:
	docker run --rm \
	--privileged \
	-v $(PWD):/go/src/$(PKG) \
	-w /go/src/$(PKG) \
	goreleaser/goreleaser \
		release \
		--rm-dist \
		--release-notes=./.changes/$(VERSION).md
		--snapshot

.PHONY: update
update:
	go get -u ./...
