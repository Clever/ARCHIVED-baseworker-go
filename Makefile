SHELL := /bin/bash
PKG = github.com/Clever/baseworker-go
PKGS = $(PKG)

.PHONY: test golint

golint:
	go get github.com/golang/lint/golint

test: $(PKGS)

README.md: *.go
	go get github.com/robertkrimen/godocdown/godocdown
	godocdown > README.md

$(PKGS): golint
	go get -d -t $@
ifneq ($(NOLINT),1)
	PATH=$(PATH):$(GOPATH)/bin golint $(GOPATH)/src/$@*/**.go
endif
ifeq ($(COVERAGE),1)
	go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	go tool cover -html=$(GOPATH)/src/$@/c.out
else
	go test $@ -test.v
endif
