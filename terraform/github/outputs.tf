output "repository_name" {
  description = "Name of the GitHub repository"
  value       = github_repository.core_course_labs.name
}

output "repository_url" {
  description = "URL of the GitHub repository"
  value       = github_repository.core_course_labs.html_url
}

output "repository_full_name" {
  description = "Full name of the GitHub repository"
  value       = github_repository.core_course_labs.full_name
}

output "default_branch" {
  description = "Default branch of the repository"
  value       = github_branch_default.default.branch
}
