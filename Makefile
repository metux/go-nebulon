# GO := go1.23.0
# GO := go1.22.0
PATH := $(PATH):/usr/lib/go-1.22/

SHELL := /bin/bash
PACKAGE := github.com/metux/go-nebulon

# GO ?= /usr/lib/go-1.22/bin/go
GO ?= go

EXECUTABLE=go-nebulon

TEST_INPUT=./go-nebulon
TEST_OUTPUT=test1.tmp

# export GOEXPERIMENT=rangefunc
#export GODEBUG=http1debug=2
#export GODEBUG=http2debug=2

run: get-deps gen-proto compile
	$(GO) version
	time ./go-nebulon
#	diff -ruN $(TEST_INPUT) $(TEST_OUTPUT)

gen-proto:
	$(MAKE) -C wire

get-deps:
	$(GO) get

test:
	@rm -Rf .storedata
	@$(GO) test -v $(PACKAGE)/... || echo " ==== self-test failed === "

compile: get-deps gen-proto
	$(GO) build $(GOTAGS) .

fmt:
	$(GO) fmt $(PACKAGE)/...

clean:
	rm -f $(EXECUTABLE)
	$(MAKE) -C wire clean
