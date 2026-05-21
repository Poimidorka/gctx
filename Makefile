BINARY := build/gctx
GOCACHE ?= $(if $(TMPDIR),$(TMPDIR),/tmp/)gctx-go-cache
ARGS ?=

.PHONY: build test run clean

build:
	@mkdir -p build
	GOCACHE="$(GOCACHE)" go build -o "$(BINARY)" .

test:
	GOCACHE="$(GOCACHE)" go test ./...

run:
	GOCACHE="$(GOCACHE)" go run . $(ARGS)

clean:
	rm -rf build
