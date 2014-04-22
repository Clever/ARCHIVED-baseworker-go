SHELL := /bin/bash
PKG = github.com/Clever/baseworker-go
SUBPKGS = worker.go
PKGS = $(PKG) $(SUBPKGS)

.PHONY: test docs

test: $(PKGS)

README.md: *.go
	go get github.com/robertkrimen/godocdown/godocdown
	godocdown > README.md

$(PKGS):
ifeq ($(LINT),1)
	golint $(GOPATH)/src/$@*/**.go
endif
	go get -d -t $@
ifeq ($(COVERAGE),1)
	go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	go tool cover -html=$(GOPATH)/src/$@/c.out
else
	go test $@ -test.v
endif
