# GO := go1.23.0
# GO := go1.22.0
# PATH := $(PATH):/usr/lib/go-1.22/

SHELL := /bin/bash
PACKAGE := github.com/metux/go-nebulon

GO ?= /usr/lib/go-1.22/bin/go

TEST_INPUT=./go-nebulon
TEST_OUTPUT=test1.tmp

#export GODEBUG=http1debug=2
#export GODEBUG=http2debug=2

run: get-deps gen-proto compile
	time ./go-nebulon
	diff -ruN $(TEST_INPUT) $(TEST_OUTPUT)

gen-proto:
	protoc -I=. --go_out=. wire/nebulon.proto --go_opt="Mwire/nebulon.proto=./wire"

get-deps:
	$(GO) get

test:
	$(GO) test -v $(PACKAGE)/...

compile: get-deps gen-proto
	$(GO) build $(GOTAGS) .

fmt:
	$(GO) fmt $(PACKAGE)/...

clean:
	rm -f eertool

#dbmigrate:
#	cd migrations && $(GO) run
