.PHONY: help
help: ## prints help (only for tasks with comment)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

EXECUTABLE_NAME=tf-summarize
COMMIT?=$(shell git describe --always 2> /dev/null)
TF_SUMMARIZE_VERSION?="development-$(COMMIT)"
build: ## build the binary
	go build -o $(EXECUTABLE_NAME) -ldflags="-X 'main.version=$(TF_SUMMARIZE_VERSION)'" .

install: build ## build and install to /usr/local/bin/
	cp $(EXECUTABLE_NAME) /usr/local/bin/$(EXECUTABLE_NAME)

test:
	go test ./...

i: install ## build and install to /usr/local/bin/
