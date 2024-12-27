# GO := go1.23.0
# GO := go1.22.0
# PATH := $(PATH):/usr/lib/go-1.22/

SHELL := /bin/bash
PACKAGE := github.com/metux/go-nebulon

GO ?= /usr/lib/go-1.22/bin/go

#export GODEBUG=http1debug=2
#export GODEBUG=http2debug=2

# GO := go

run: get-deps gen-proto compile
#	time $(GO) run $(GOTAGS) .
	time ./go-nebulon
	diff -ruN go-nebulon go-nebulon.tmp

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
