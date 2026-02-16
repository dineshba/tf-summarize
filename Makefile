NAME:=tf-summarize
TERRAFORM_VERSION:=$(shell cat example/.terraform-version)
GORELEASER=go run github.com/goreleaser/goreleaser/v2@v2.13.3
GOLANGCI_LINT=go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.0
GOSEC=go run github.com/securego/gosec/v2/cmd/gosec@v2.23.0
VERSION:=$(shell cat VERSION)

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

default: build

.PHONY: help
help: ## prints help (only for tasks with comment)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

version: # print version
	@echo $(VERSION)
.PHONY: version

# TODO: dynamically set architecture, which is currently hard-coded to amd64
install: build ## build and install to /usr/local/bin/
	cp dist/$(NAME)_$(shell echo $(shell uname) | tr '[:upper:]' '[:lower:]')_amd64*/$(NAME) /usr/local/bin/$(NAME)

test: lint ## go test
	go test -v ./... -count=1

i: install ## build and install to /usr/local/bin/

lint: ## lint source code
	$(GOLANGCI_LINT) run --timeout 10m -v
.PHONY: lint

gosec: ## run gosec security scanner
	$(GOSEC) -exclude=G204,G705 ./...
.PHONY: lint

example: ## generate example Terraform plan
	$(call generate-example,$(TERRAFORM_VERSION))
.PHONY: example

build: ## build and test
	$(GORELEASER) release \
	--snapshot \
	--skip=publish,sign \
	--clean
.PHONY: build

tag: ## create $(VERSION) git tag
	echo "creating git tag $(VERSION)"
	git tag $(VERSION)
	git push origin $(VERSION)
.PHONY: tag

release: ## release $(VERSION)
	$(GORELEASER) release \
		--clean
.PHONY: release
