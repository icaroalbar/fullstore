.PHONY: fmt vet lint test check run web-dev web-test web-lint web-typecheck web-check

export PATH := $(shell go env GOPATH)/bin:$(PATH)

# Enquanto o repo não tem arquivos .go, vet, lint e test não têm o que checar.
HAS_GO := $(shell find . -name '*.go' -not -path './.git/*' -not -path './web/*' -print -quit)

fmt:
	@out="$$(gofmt -l .)"; \
	if [ -n "$$out" ]; then echo "gofmt: arquivos fora do formato (rode gofmt -w):"; echo "$$out"; exit 1; fi

vet:
ifeq ($(HAS_GO),)
	@echo "vet: nenhum arquivo Go ainda"
else
	go vet ./...
endif

lint:
ifeq ($(HAS_GO),)
	@echo "lint: nenhum arquivo Go ainda"
else
	@if [ -d web/node_modules ]; then mv web/node_modules web/.node_modules.bak; fi; \
	golangci-lint run; \
	RESULT=$$?; \
	if [ -d web/.node_modules.bak ]; then mv web/.node_modules.bak web/node_modules; fi; \
	exit $$RESULT
endif

test:
ifeq ($(HAS_GO),)
	@echo "test: nenhum arquivo Go ainda"
else
	go test -race -cover ./...
endif

check: fmt vet lint test

run:
	go run ./cmd/api

web-dev:
	npm --prefix web run dev

web-test:
	npm --prefix web test -- --run

web-check: web-lint web-typecheck web-test web-build

web-lint:
	npm --prefix web run lint

web-typecheck:
	npm --prefix web run typecheck

web-build:
	npm --prefix web run build
