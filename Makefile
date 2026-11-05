.PHONY: fmt vet lint test check run

export PATH := $(shell go env GOPATH)/bin:$(PATH)

# Enquanto o repo não tem arquivos .go, vet, lint e test não têm o que checar.
HAS_GO := $(shell find . -name '*.go' -not -path './.git/*' -print -quit)

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
	golangci-lint run
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
