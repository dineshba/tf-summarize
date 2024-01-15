TERRAFORM_VERSION:=$(shell cat example/.terraform-version)

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

test: lint
	go test ./...

i: install ## build and install to /usr/local/bin/

lint:
	golangci-lint run --timeout 10m -v

define generate-example
	docker run \
		--interactive \
		--tty \
		--volume $(shell pwd):/src \
		--workdir /src/example \
		--entrypoint /bin/sh \
		hashicorp/terraform:$(1) \
			-c \
				"terraform init && \
				terraform plan -out tfplan && \
				terraform show -json tfplan > tfplan.json"
endef

example:
	$(call generate-example,$(TERRAFORM_VERSION))
.PHONY: example
