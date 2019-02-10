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
	  github.com/Songmu/ghch/cmd/ghch

test: deps
	go test

lint: devel-deps
	go vet
	golint -set_exit_status

cover: devel-deps
	goveralls

bump: devel-deps
	_tools/releng

release: bump

.PHONY: test deps devel-deps lint cover build bump release
