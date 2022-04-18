## Terraform Plan Summary

[![Build](https://github.com/dineshba/terraform-plan-summary/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/dineshba/terraform-plan-summary/actions/workflows/build.yml) [![goreleaser](https://github.com/dineshba/terraform-plan-summary/actions/workflows/release.yml/badge.svg)](https://github.com/dineshba/terraform-plan-summary/actions/workflows/release.yml) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/dineshba/terraform-plan-summary)

`tf-plan-summarize` is a command-line utility to print the summary of the terraform plan

### Why do we need it ?

##### For below use-cases:
- Most of the time, we make changes to the terraform files or tf-var files and run the plan command. And we precisly know which resources will get affected. In those time, we would like to just see the resource name and it's change.
- When our plan have more than say 10 changes, we will first see what are the deleted changes or we will just see the list of resources that get affected.

![demo](example/demo.gif)

### Install

#### Download zip in release page
1. Go to release page [https://github.com/dineshba/terraform-plan-summary/releases](https://github.com/dineshba/terraform-plan-summary/releases)
2. Download the zip for your OS and unzip it
3. Copy it to local bin using `cp tf-plan-summarize /usr/local/bin/tf-plan-summarize` or to location which is part of `$PATH`
4. (For Mac Only) Give access to run if prompted. [Refer here](https://stackoverflow.com/a/19551359/5305962)

#### Clone and Build Binary
1. Clone this repo
2. Build binary using `make build` or `go build -o tf-plan-summarize .`
3. Install it to local bin using `make install` or `cp tf-plan-summarize /usr/local/bin/tf-plan-summarize`

### Usage

```sh
$ tf-plan-summarize -h

Usage of tf-plan-summarize [args] [tf-plan.json|tfplan]

  -draw
        [Optional, used only with -tree or -separate-tree] draw trees instead of plain tree
  -out string
        [Optional] write output to file
  -separate-tree
        [Optional] print changes in tree format for each add/delete/change/recreate changes
  -tree
        [Optional] print changes in tree format
```

### Example

#### Simple Example
```sh
# run terraform plan command
terraform plan -out=tfplan
# provide json output from plan
terraform show -json tfplan | tf-plan-summarize # will print the summary in stdout in table format
# provide plan output directly
tf-plan-summarize tfplan
```

#### More Examples
```sh
# terraform plan as input based examples
terraform plan -out=tfplan
tf-plan-summarize tfplan                           # summary in table format
tf-plan-summarize -tree tfplan                     # summary in tree format
tf-plan-summarize -tree -draw tfplan               # summary in 2D tree format
tf-plan-summarize -separate-tree tfplan            # summary in separate tree format
tf-plan-summarize -separate-tree -draw tfplan      # summary in separate 2D tree format
tf-plan-summarize -out=summary.md tfplan           # summary in output file instead of stdout

# json as input based examples
terraform plan -out=tfplan
terraform show -json tfplan > output.json
tf-plan-summarize output.json                                 # summary in table format
cat output.json | tf-plan-summarize                           # summary in table format
cat output.json | tf-plan-summarize -tree                     # summary in tree format
cat output.json | tf-plan-summarize -tree -draw               # summary in 2D tree format
cat output.json | tf-plan-summarize -separate-tree            # summary in separate tree format
cat output.json | tf-plan-summarize -separate-tree -draw      # summary in separate 2D tree format
cat output.json | tf-plan-summarize -out=summary.md           # summary in output file instead of stdout

```

### TODO

- [x] Read terraform state file directly. (Currently need to convert to json and pass it)
- [ ] Directly run the terraform plan and show the summary
- [ ] Able to show summary of the current terraform state
- [ ] Include version subcommand in binary
