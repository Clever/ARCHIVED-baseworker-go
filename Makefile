SHELL := /bin/bash
PKG = github.com/Clever/baseworker-go

.PHONY: test docs

run:
	GEARMAN_HOST=localhost GEARMAN_PORT=4730 go run main.go

docs:
	go get github.com/robertkrimen/godocdown/godocdown
	godocdown > README.md

$(PKG):
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
