.PHONY: help
help: ## prints help (only for tasks with comment)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


build: ## build the binary
	go build terraform-plan-summary.go

install: build ## build and install to /usr/local/bin/
	cp terraform-plan-summary /usr/local/bin/terraform-plan-summary

i: install ## build and install to /usr/local/bin/