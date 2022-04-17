resource "github_repository" "repository" {
  name        = var.name
  description = var.description
  topics      = var.topics

  visibility = "public"

  has_downloads        = true
  has_issues           = true
  has_projects         = true
  has_wiki             = true
  vulnerability_alerts = false
}

resource "github_branch" "main" {
  repository = github_repository.repository.name
  branch     = "main"
}

resource "github_branch" "development" {
  repository = github_repository.repository.name
  branch     = "development"
}