SHELL := /bin/bash
PKG = github.com/Clever/baseworker-go

.PHONY: test

run:
	GEARMAN_HOST=localhost GEARMAN_PORT=4730 go run main.go

$(PKG):
	go get -d -t $@
	go test $@ -test.v
