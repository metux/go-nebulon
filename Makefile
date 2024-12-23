# GO := go1.23.0
# GO := go1.22.0
# PATH := $(PATH):/usr/lib/go-1.22/

PACKAGE := github.com/metux/go-nebulon

GO ?= /usr/lib/go-1.22/bin/go

#export GODEBUG=http1debug=2
#export GODEBUG=http2debug=2

# GO := go

run:
	$(GO) get
	$(GO) run $(GOTAGS) .

test:
	$(GO) test -v $(PACKAGE)/...

compile:
	$(GO) get
	$(GO) build $(GOTAGS) .

fmt:
	$(GO) fmt $(PACKAGE)/...

clean:
	rm -f eertool

#dbmigrate:
#	cd migrations && $(GO) run
