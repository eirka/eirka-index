# Verify targets. Hermes/Claude workers pick up build/lint/test/check from this
# file; keep the target names. See CLAUDE.md. `build` uses ./... so no binary
# lands in the tree (nothing is gitignored here).

GO ?= go

.PHONY: build lint test check

build:
	$(GO) build ./...

lint:
	@out="$$(gofmt -l -s $$($(GO) list -f '{{.Dir}}' ./...))"; if [ -n "$$out" ]; then echo "gofmt -s -w needed:"; echo "$$out"; exit 1; fi
	$(GO) vet ./...

test:
	$(GO) test -count=1 ./...

check: build lint test
