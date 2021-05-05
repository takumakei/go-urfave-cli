.DEFAULT_GOAL := help

.PHONY: help snapshot

help: ## show this message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-8s\033[0m %s\n", $$1, $$2}'

snapshot: ## goreleaser build --rm-dist --snapshot
	goreleaser build --rm-dist --snapshot
