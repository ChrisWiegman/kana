PKG       := github.com/ChrisWiegman/kana-cli
VERSION   := $(shell git describe --tags || echo "0.0.1")
GITHASH   := $(shell git describe --tags --always --dirty)
TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
ARGS = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

%:
	@:

.PHONY: build
build:
	go mod vendor
	go run ./scripts/generateTemplateConstants.go
	go install \
		-ldflags "-s -w -X $(PKG)/internal/cmd.Version=$(VERSION) -X $(PKG)/internal/cmd.GitHash=$(GITHASH) -X $(PKG)/internal/cmd.Timestamp=$(TIMESTAMP)" \
		./cmd/...

.PHONY: change
change:
	docker run \
		--rm \
		--platform linux/amd64 \
		--mount type=bind,source=$(PWD),target=/src \
		-w /src \
		-it \
		ghcr.io/miniscruff/changie \
		new

.PHONY: changelog
changelog:
	docker run \
		--rm \
		--platform linux/amd64 \
		--mount type=bind,source=$(PWD),target=/src \
		-w /src \
		-it \
		ghcr.io/miniscruff/changie \
		batch $(call ARGS,defaultstring)
	docker run \
		--rm \
		--platform linux/amd64 \
		--mount type=bind,source=$(PWD),target=/src \
		-w /src \
		-it \
		ghcr.io/miniscruff/changie \
		merge

.PHONY: clean
clean:
	rm -rf \
		dist \
		vendor \
		internal/config/constants.go

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

.PHONY: release
release:
	echo $(VERSION)
	docker run --rm \
	--privileged \
	-v $(PWD):/go/src/$(PKG) \
	-w /go/src/$(PKG) \
	goreleaser/goreleaser \
		release \
		--rm-dist \
		--release-notes=./.changes/$(VERSION).md \
		--snapshot
