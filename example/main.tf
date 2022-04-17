terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "4.23.0"
    }
  }
}

provider "github" {
  owner = "dineshba"
}

locals {
  repos = [
    {
      name        = "terraform-plan-summary"
      description = "A command-line utility to print the summary of the terraform plan"
      topics = [
        "summary",
        "terraform",
      ]
    },
    {
      name        = "demo-repository"
      description = "description of the repo"
      topics = [
        "tags",
      ]
    }
  ]
}

module "github" {
  source = "./github"
  for_each = {
    for repo in local.repos : repo.name => repo
  }

  name        = each.key
  description = each.value.description
  topics      = each.value.topics
}

resource "github_repository" "terraform_plan_summary" {
  name        = "terraform-plan-summary"
  description = "A command-line utility to print the summary of the terraform plan"
  topics = [
    "summary",
    "terraform",
  ]

  visibility = "public"

  has_downloads        = true
  has_issues           = true
  has_projects         = true
  has_wiki             = true
  vulnerability_alerts = false
}
