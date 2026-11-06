VERSION ?= $(shell git describe --tags --always)
BUILD_LDFLAGS = "-s -w"
u := $(if $(update),-u)

.PHONY: deps
deps:
	go get ${u}
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	go run honnef.co/go/tools/cmd/staticcheck@v0.8.2 ./...

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/timestamp

.PHONY: install
install: deps
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/timestamp

.PHONY: crossbuild
crossbuild:
	go run github.com/Songmu/goxz/cmd/goxz@v0.13.0 -pv=$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist/$(VERSION) ./cmd/timestamp

.PHONY: upload
upload:
	gh release create $(VERSION) dist/$(VERSION)/* --verify-tag --generate-notes
