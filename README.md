## tf-summarize (Terraform Summarizer)

[![Build](https://github.com/dineshba/terraform-plan-summary/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/dineshba/terraform-plan-summary/actions/workflows/build.yml)
[![goreleaser](https://github.com/dineshba/terraform-plan-summary/actions/workflows/release.yml/badge.svg)](https://github.com/dineshba/terraform-plan-summary/actions/workflows/release.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/dineshba/terraform-plan-summary)
[![](https://img.shields.io/badge/License-MIT%202.0-blue.svg)](https://github.com/dineshba/tf-summarize/blob/main/LICENSE)

`tf-summarize` is a command-line utility to print the summary of the terraform plan

<!-- TOC -->
  * [tf-summarize (Terraform Summarizer)](#tf-summarize-terraform-summarizer)
    * [Demo](#demo)
    * [Why do we need it ?](#why-do-we-need-it-)
    * [Install](#install)
      * [Using Go](#using-go)
      * [Using Brew](#using-brew)
      * [Using asdf](#using-asdf)
      * [Using Docker](#using-docker)
      * [Using Github action](#using-github-action)
      * [Download zip in release page](#download-zip-in-release-page)
      * [Clone and Build Binary](#clone-and-build-binary)
    * [Usage](#usage)
    * [Examples](#examples)
      * [GitHub Actions Workflow](#github-actions-workflow)
        * [Using setup-tf-summarize action](#using-setup-tf-summarize-action)
        * [Other approach](#other-approach)
      * [Codefresh example](#codefresh-example)
      * [Comment terraform plan summary in PRs](#comment-terraform-plan-summary-in-prs)
      * [Interactive summary review](#interactive-summary-review)
    * [Screenshot](#screenshot)
    * [TODO](#todo)
<!-- TOC -->
### Demo

![demo](example/demo.gif)

If demo is slower for you, please see [examples](#examples) and [screenshot](#screenshot) for more details

### Why do we need it ?

- Most of the time, we make changes to the terraform files or tf-var files and run the plan command. And we precisely know which resources will get affected. In those time, we would like to just see the resource name and it's change.
- When our plan have more than say 10 changes, we will first see what are the deleted changes, or we will just see the list of resources that get affected.

### Install

#### Using Go
```sh
go install github.com/dineshba/tf-summarize@latest
```

#### Using Brew
```sh
brew tap thecasualcoder/stable
brew install tf-summarize
```

#### Using asdf
```sh
asdf plugin add tf-summarize
asdf install tf-summarize latest
```

> See the asdf [plugin source](https://github.com/adamcrews/asdf-tf-summarize#install) for more information.

#### Using Docker
```sh
docker run -v $PWD:/workspace -w /workspace ghcr.io/dineshba/tf-summarize -v # prints version
docker run -v $PWD:/workspace -w /workspace ghcr.io/dineshba/tf-summarize tfplan.json
docker run -v $PWD:/workspace -w /workspace ghcr.io/dineshba/tf-summarize -tree tfplan.json
# see example section for more examples


# add alias if needed
alias tf-summarize="docker run -v $PWD:/workspace -w /workspace ghcr.io/dineshba/tf-summarize"
tf-summarize tfplan.json
# see example section for more examples
```

> tf-summarize will accept tfplan directly. Docker based tf-summarize is not having terraform-binary to convert tfplan to json. So it only works with json files.
> If you want to use with tfplan directly, You can build a base image with terraform (use same version as tfplan created) and have our binary in it.

#### Using Github action

The below example uses [kishaningithub/setup-tf-summarize](https://github.com/kishaningithub/setup-tf-summarize)

```yaml
steps:
- uses: kishaningithub/setup-tf-summarize@v1
```

For more customization options, kindly refer to the [documentation here](https://github.com/kishaningithub/setup-tf-summarize#usage)

#### Download zip in release page
1. Go to release page [https://github.com/dineshba/terraform-plan-summary/releases](https://github.com/dineshba/terraform-plan-summary/releases)
2. Download the zip for your OS and unzip it
3. Copy it to local bin using `cp tf-summarize /usr/local/bin/tf-summarize` or to location which is part of `$PATH`
4. (For Mac Only) Give access to run if prompted. [Refer here](https://stackoverflow.com/a/19551359/5305962)

#### Clone and Build Binary
1. Clone this repo
2. Build binary using `make build` or `go build -o tf-summarize .`
3. Install it to local bin using `make install` or `cp tf-summarize /usr/local/bin/tf-summarize`

### Usage

```sh
$ tf-summarize -h

Usage of tf-summarize [args] [tf-plan.json|tfplan]

  -draw
        [Optional, used only with -tree or -separate-tree] draw trees instead of plain tree
  -html
        [Optional] print changes in html format
  -json
        [Optional] print changes in json format
  -json-sum
        [Optional] print summary in json format
  -md
        [Optional, used only with table view] output table as markdown
  -out string
        [Optional] write output to file
  -separate-tree
        [Optional] print changes in tree format for add/delete/change/recreate changes
  -tree
        [Optional] print changes in tree format
  -v    print version
```

### Examples

```sh
# run terraform plan command
terraform plan -out=tfplan

# provide plan itself directly
tf-summarize tfplan                           # summary in table format
tf-summarize -tree tfplan                     # summary in tree format
tf-summarize -tree -draw tfplan               # summary in 2D tree format
tf-summarize -json tfplan                     # summary in json format
tf-summarize -json-sum tfplan                 # summary of changes in json format
tf-summarize -separate-tree tfplan            # summary in separate tree format
tf-summarize -separate-tree -draw tfplan      # summary in separate 2D tree format
tf-summarize -out=summary.md tfplan           # summary in output file instead of stdout

# provide json output from plan
terraform show -json tfplan | tf-summarize    # summary in table format
terraform show -json tfplan > output.json
tf-summarize output.json                      # summary in table format
```

#### GitHub Actions Workflow

##### Using setup-tf-summarize action

The below example uses [kishaningithub/setup-tf-summarize](https://github.com/kishaningithub/setup-tf-summarize)

```yaml
name: Run tf-summarize

on:
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: kishaningithub/setup-tf-summarize@v1

      - name: Print tf-summarize version and help
        run: |
          tf-summarize -v
          tf-summarize -h
```

For more customization options, kindly refer to the [documentation here](https://github.com/kishaningithub/setup-tf-summarize#usage)

##### Other approach

Please refer this sample [GitHub actions file](.github/workflows/demo.yml) and the sample runs [here](https://github.com/dineshba/tf-summarize/actions/workflows/demo.yml)

> Note: If you are using `hashicorp/setup-terraform` GitHub action to set up terraform, ensure terraform_wrapper is set to false.

#### Codefresh example

The below assumes that you have a stage called `show` in your `codefresh.yaml` example

We need to output the plan to json as there were some strange issues on codefresh when we used the `tfplan` file

A full example is available at [`example/codefresh/codefresh.yaml`](/example/codefresh/codefresh.yaml)

```yaml
  TerraformPlan:
    title: Terraform Plan
    image: hashicorp/terraform:light
    stage: plan
    working_directory: "${{clone}}"
    commands:
      - terraform plan -out=tfplan
      - terraform show -json tfplan > output.json
  
  tfSummarize:
    title: Show Changes
    image: ghcr.io/dineshba/tf-summarize
    stage: show
    working_directory: "${{clone}}"
    commands:
      - tf-summarize output.json
```

> :warning: **Not maintained by [dineshba](https://github.com/dineshba)** : The above example is maintained by [userbradley](https://github.com/userbradley) - Any questions related please tag `@userbradley` on an issue 

#### Comment terraform plan summary in PRs

Refer [this example](https://github.com/dineshba/tf-summarize/blob/demo-pr/.github/workflows/demo.yml#L61-L73) to add comments in your PR. Sample [comment](https://github.com/dineshba/tf-summarize/pull/19#issuecomment-1295882938) added by github actions bot.

#### Interactive summary review

You can use tool [fx](https://github.com/antonmedv/fx) to review the summary fo the terraform change
```sh
tf-summarize -json tfplan | fx
```

### Screenshot
![screenshot](example/tf-summarize-screenshot.png)

### TODO

- [x] Read terraform state file directly. (Currently need to convert to json and pass it)
- [ ] Directly run the terraform plan and show the summary
- [ ] Able to show summary of the current terraform state
- [x] Include version subcommand in binary
