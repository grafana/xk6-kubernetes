MAKEFLAGS += --silent
WORKFLOW        ?= .github/workflows/k6-ci.yml
K6_CI_REF       := $(shell grep -oE 'grafana/k6-ci/[^@[:space:]]+@[A-Za-z0-9._/-]+' $(WORKFLOW) | head -n1 | cut -d@ -f2)
LINT_BASE       ?= .golangci-base.yml
LINT_CONFIG     ?= .golangci.yml
LINT_PATCH      ?= .golangci.patch
LINT_CONFIG_URL := https://raw.githubusercontent.com/grafana/k6-ci/$(K6_CI_REF)/.golangci.yml

all: clean format lint test build

## help: Prints a list of available build targets.
help:
	echo "Usage: make <OPTIONS> ... <TARGETS>"
	echo ""
	echo "Available targets are:"
	echo ''
	sed -n 's/^##//p' ${PWD}/Makefile | column -t -s ':' | sed -e 's/^/ /'
	echo
	echo "Targets run by default are: `sed -n 's/^all: //p' ./Makefile | sed -e 's/ /, /g' | sed -e 's/\(.*\), /\1, and /'`"

## clean: Removes any previously created build artifacts.
clean:
	rm -f ./k6 $(LINT_BASE) $(LINT_CONFIG)

## build: Builds a custom 'k6' with the local extension. 
build:
	go install go.k6.io/xk6/cmd/xk6@latest
	xk6 build --with $(shell go list -m)=.

## format: Applies Go formatting to code.
format:
	go fmt ./...

$(LINT_BASE): $(WORKFLOW)
	curl -fsSL $(LINT_CONFIG_URL) -o $@

$(LINT_CONFIG): $(LINT_BASE) $(LINT_PATCH)
	cp $(LINT_BASE) $(LINT_CONFIG)
	git apply $(LINT_PATCH)

## lint: Runs golangci-lint with the patched k6-ci configuration.
lint: $(LINT_CONFIG)
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$$(head -n1 $(LINT_BASE) | tr -d '# ') \
	  run --config=$(LINT_CONFIG) ./...

## test: Executes any unit tests.
test:
	go test -cover -race ./...

.PHONY: build clean format help lint test
