VERSION = $(shell godzil show-version)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/Songmu/timestamper.revision=$(CURRENT_REVISION)"
u := $(if $(update),-u)

.PHONY: deps
deps:
	go get ${u}
	go mod tidy

.PHONY: devel-deps
devel-deps:
	go install honnef.co/go/tools/cmd/staticcheck@v0.8.2
	go install github.com/tcnksm/ghr@v0.18.6
	go install github.com/Songmu/godzil/cmd/godzil@v0.23.2

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint: devel-deps
	staticcheck ./...

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/timestamp

.PHONY: install
install: deps
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/timestamp

.PHONY: release
release:
	godzil release

CREDITS: deps devel-deps go.sum
	godzil credits -w

.PHONY: crossbuild
crossbuild: devel-deps
	godzil crossbuild -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist/v$(VERSION) ./cmd/timestamp

.PHONY: upload
upload:
	ghr -body="$$(godzil changelog --latest -F markdown)" v$(VERSION) dist/v$(VERSION)
