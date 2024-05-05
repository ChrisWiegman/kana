PKG          := github.com/ChrisWiegman/kana
VERSION      := $(shell git describe --tags || echo "0.0.1")
TIMESTAMP    := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
ARGS          = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

%:
	@:

.PHONY: build
build:
	go mod vendor
	go$(GO_VERSION) build \
		-o ./build/kana \
		-ldflags "-s -w -X $(PKG)/internal/cmd.Version=$(VERSION) -X $(PKG)/internal/cmd.Timestamp=$(TIMESTAMP)" \
		./cmd/...

.PHONY: build-test-image
build-test-image:
	docker build -t kana-test .

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
		build

.PHONY: install
install:
	go mod vendor
	go$(GO_VERSION) install \
		-ldflags "-s -w -X $(PKG)/internal/cmd.Version=$(VERSION) -X $(PKG)/internal/cmd.Timestamp=$(TIMESTAMP)" \
		./cmd/...

.PHONY: lint
lint:
	@if [ ! -f $GOPATH/bin/gilangci-lint  ]; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest;\
	fi
	@golangci-lint \
			run

.PHONY: mockery
mockery:
	docker \
		run \
		--rm \
		--mount type=bind,source=$(PWD),target=/src \
		-w /src/internal/docker \
		vektra/mockery \
		--all

.PHONY: update
update:
	go get -u ./...
	go mod vendor
	go mod tidy

.PHONY: snapshot
snapshot:
	docker run --rm \
	--privileged \
	-v $(PWD):/go/src/$(PKG) \
	-w /go/src/$(PKG) \
	goreleaser/goreleaser \
		release \
		--clean \
		--release-notes=./.changes/$(VERSION).md \
		--snapshot

.PHONY: test
test: clean build-test-image
	docker run --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):/usr/src/kana \
		-w /usr/src/kana \
		kana-test
	$(MAKE) clean
	$(MAKE) install

.PHONY: update-test-snapshot
update-test-snapshot:
	go build \
        -o ./build/kana \
        -buildvcs=false \
        -ldflags "-s -w -X github.com/ChrisWiegman/kana/internal/cmd.Version=1.0.0 -X github.com/ChrisWiegman/kana/internal/cmd.Timestamp=2024-03-16_10:50:11PM" \
        ./cmd/... && \
    UPDATE_SNAPS=true go test \
        -v \
        -timeout 30s\
        -cover \
        ./...