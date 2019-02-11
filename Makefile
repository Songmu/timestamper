VERSION = $(shell gobump show -r)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/Songmu/timestamper.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

deps:
	go get ${u} -d -t ./...

test-deps:
	go get ${u} -d -v -t ./...

devel-deps: deps
	go get ${u} golang.org/x/lint/golint   \
	  github.com/mattn/goveralls           \
	  github.com/motemen/gobump/cmd/gobump \
	  github.com/Songmu/goxz/cmd/goxz      \
	  github.com/Songmu/ghch/cmd/ghch      \
	  github.com/tcnksm/ghr

test: deps
	go test

lint: devel-deps
	go vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/timestamp

bump: devel-deps
	_tools/releng

crossbuild:
	goxz -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
	  -os=linux,darwin,windows -d=./dist/v$(VERSION) ./cmd/timestamp

upload:
	ghr v$(VERSION) dist/v$(VERSION)

release: bump crossbuild upload

.PHONY: test deps devel-deps lint cover build bump release
