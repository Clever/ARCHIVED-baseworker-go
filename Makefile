SHELL := /bin/bash
PKG = github.com/Clever/baseworker-go
PKGS = $(PKG)

.PHONY: test golint README

golint:
	@go get github.com/golang/lint/golint

test: $(PKGS)

README.md: *.go
	@go get github.com/robertkrimen/godocdown/godocdown
	$(GOPATH)/bin/godocdown $(PKG) > README.md
README: README.md

$(PKGS): golint README
	@if [[ -z "$(DRONE)" ]]; then \
		echo "resetting gearman" \
		./reset_gearmand.sh; \
	fi
	@go get -d -t $@
	@gofmt -w=true $(GOPATH)/src/$@*/**.go
ifneq ($(NOLINT),1)
	@echo "LINTING..."
	@PATH=$(PATH):$(GOPATH)/bin golint $(GOPATH)/src/$@*/**.go
	@echo ""
endif
ifeq ($(COVERAGE),1)
	@go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	@go tool cover -html=$(GOPATH)/src/$@/c.out
else
	@echo "TESTING..."
	@go test $@ -test.v
endif
