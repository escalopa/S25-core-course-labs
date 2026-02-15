variable "github_token" {
  description = "GitHub personal access token (use GITHUB_TOKEN env var instead)"
  type        = string
  sensitive   = true
}

variable "github_owner" {
  description = "GitHub username or organization"
  type        = string
  default     = "escalopa"
}

variable "repo_name" {
  description = "Name of the GitHub repository"
  type        = string
  default     = "S25-core-course-labs"
}

variable "repo_description" {
  description = "Description of the GitHub repository"
  type        = string
  default     = "S25 Core Course Labs - DevOps Engineering"
}

variable "repo_visibility" {
  description = "Visibility of the repository (public or private)"
  type        = string
  default     = "public"
}

variable "default_branch" {
  description = "Default branch of the repository"
  type        = string
  default     = "master"
}
