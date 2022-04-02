GO ?= go

.PHONY: test

# Run test
test:
	$(GO) test ./...
