terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
  }

  required_version = ">= 1.0"
}

# Token is read from the TF_VAR_github_token environment variable
provider "github" {
  token = var.github_token
}

resource "github_repository" "core_course_labs" {
  name        = var.repo_name
  description = var.repo_description

  visibility = var.repo_visibility

  has_issues   = true
  has_wiki     = true
  has_projects = true

  auto_init          = false
  allow_merge_commit = true
  allow_squash_merge = true
  allow_rebase_merge = true
}

resource "github_branch_default" "default" {
  repository = github_repository.core_course_labs.name
  branch     = var.default_branch
}

resource "github_branch_protection" "default" {
  repository_id = github_repository.core_course_labs.node_id

  pattern          = var.default_branch
  enforce_admins   = true
  allows_deletions = false

  required_pull_request_reviews {
    dismiss_stale_reviews      = true
    required_approving_review_count = 1
  }
}
