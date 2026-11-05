.DEFAULT_GOAL := help

BINARY := inkcheck
TAG := $(shell git describe --tags --exact-match HEAD 2>/dev/null)

.PHONY: help build install test vet fmt check snapshot release clean

help: ## Show this help
	@echo "Usage: make <target>"
	@echo
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the CLI into ./inkcheck
	go build -o $(BINARY) ./cmd/inkcheck

install: ## Install the CLI into $GOPATH/bin
	go install ./cmd/inkcheck

test: ## Run the tests
	go test ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format the code with gofmt
	gofmt -w .

check: vet test ## Run vet and tests

snapshot: check ## Build release archives locally without signing or publishing
	goreleaser release --snapshot --clean

release: check ## Tag HEAD as TAG=vX.Y.Z (if needed), push, then sign, notarise and publish
	@test -z "$$(git status --porcelain)" || { echo "Working tree is not clean."; exit 1; }
	@test "$$(git rev-parse --abbrev-ref HEAD)" = main || { echo "Not on main."; exit 1; }
	@test -n "$(TAG)" || { echo "HEAD has no tag. Run: make release TAG=vX.Y.Z"; exit 1; }
	@echo "$(TAG)" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$$' || { echo "TAG must look like vX.Y.Z, got '$(TAG)'."; exit 1; }
	@if git rev-parse -q --verify "refs/tags/$(TAG)" >/dev/null; then \
		test "$$(git rev-list -n 1 $(TAG))" = "$$(git rev-parse HEAD)" || { echo "Tag $(TAG) already exists and is not on HEAD."; exit 1; }; \
	else \
		echo "git tag -a $(TAG) -m '$(TAG)'"; \
		git tag -a $(TAG) -m '$(TAG)'; \
	fi
	git push origin main $(TAG)
	GORELEASER_CURRENT_TAG=$(TAG) GITHUB_TOKEN=$${GITHUB_TOKEN:-$$(gh auth token)} goreleaser release --clean

clean: ## Remove build output
	rm -rf $(BINARY) dist/
