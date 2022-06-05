.PHONY: build
build:
	go build ./cmd/...

.PHONY: run
run:
	go run ./cmd/...

.PHONY: install
install:
	go install ./cmd/...
