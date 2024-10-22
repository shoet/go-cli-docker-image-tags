.DEFAULT_GOAL := help

.PHONY: build clean

build: ## build binary file
	go build -ldflags="-s -w" -o ./.bin/docker-image-tags

clean: ## clean up
	rm -rf ./.bin

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
