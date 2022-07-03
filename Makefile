.PHONY: build
build:
	go mod vendor
	go run ./scripts/generateTemplateConstants.go
	go build ./cmd/...

.PHONY: run
run:
	go mod vendor
	go run ./scripts/generateTemplateConstants.go
	go run ./cmd/...

.PHONY: install
install:
	go mod vendor
	go run ./scripts/generateTemplateConstants.go
	go install ./cmd/...
