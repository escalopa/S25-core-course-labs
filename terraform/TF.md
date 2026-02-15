# Terraform Lab Documentation

## Table of Contents

- [Task 1: Docker Infrastructure](#task-1-docker-infrastructure)
- [Task 1: Yandex Cloud Infrastructure](#task-1-yandex-cloud-infrastructure)
- [Task 2: GitHub Infrastructure](#task-2-github-infrastructure)
- [Best Practices](#best-practices)

---

## Task 1: Docker Infrastructure

### Setup

```bash
cd terraform/docker
terraform init
terraform plan
terraform apply
```

### Terraform State

#### `terraform state list`

```
docker_container.nginx
docker_image.nginx
```

#### `terraform state show docker_container.nginx`

```
# docker_container.nginx:
resource "docker_container" "nginx" {
    name  = "lab4-nginx"
    image = docker_image.nginx.image_id
    ports {
        internal = 80
        external = 8080
    }
}
```

### Applied Changes Log

```diff
Terraform will perform the following actions:

  # docker_image.nginx will be created
  + resource "docker_image" "nginx" {
      + id          = (known after apply)
      + image_id    = (known after apply)
      + name        = "nginx:latest"
      + keep_locally = false
    }

  # docker_container.nginx will be created
  + resource "docker_container" "nginx" {
      + name  = "lab4-nginx"
      + image = (known after apply)
      + ports {
          + internal = 80
          + external = 8080
        }
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.
```

### Using Input Variables

Renamed the container using input variables:

```bash
terraform apply -var 'container_name=my-custom-nginx'
```

### Terraform Output

```
$ terraform output
container_id = "<container-id>"
container_name = "lab4-nginx"
image_id = "sha256:<hash>"
```

---

## Task 1: Yandex Cloud Infrastructure

### Setup

1. Created a Yandex Cloud account at [cloud.yandex.com](https://cloud.yandex.com/)
2. Installed the Yandex Cloud CLI (`yc`)
3. Initialized the CLI with `yc init`
4. Created an OAuth token for Terraform authentication

### Configuration

1. Created `terraform/yandex-cloud/` directory with:
   - `main.tf` — Provider config, VM instance, VPC network, and subnet
   - `variables.tf` — Input variables (token, cloud ID, folder ID, zone, VM name, image ID, SSH key path)
   - `outputs.tf` — Output values (VM ID, name, external/internal IPs, subnet ID)

2. Copied `terraform.tfvars.example` to `terraform.tfvars` and filled in actual values

### Steps to Deploy

```bash
cd terraform/yandex-cloud
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your actual Yandex Cloud credentials

terraform init
terraform plan
terraform apply
```

### Resources Created

- **VM Instance**: `lab4-vm` — 2 cores, 2GB RAM, 5% core fraction (free-tier eligible)
- **VPC Network**: `lab4-network`
- **VPC Subnet**: `lab4-subnet` — CIDR `10.2.0.0/16` in `ru-central1-a`

### Challenges

- Required a VPN to access Yandex Cloud services from certain regions
- The free-tier VM uses `core_fraction = 5` which limits CPU usage to 5% guaranteed (burstable)
- Image IDs can change over time; verified the current Ubuntu 22.04 LTS image ID via `yc compute image list --folder-id standard-images`

---

## Task 2: GitHub Infrastructure

### Setup

1. Created a GitHub Personal Access Token with `repo` and `admin:org` permissions
2. Set the token as an environment variable:

```bash
export TF_VAR_github_token="your-github-token"
```

### Import Existing Repository

Instead of creating a new repository, imported the existing one:

```bash
cd terraform/github
terraform init
terraform import "github_repository.core_course_labs" "S25-core-course-labs"
terraform import "github_branch_default.default" "S25-core-course-labs"
```

### Apply Configuration

```bash
terraform plan
terraform apply
```

### Resources Managed

- **Repository**: `S25-core-course-labs` — public, with issues/wiki/projects enabled
- **Default Branch**: `master`
- **Branch Protection**: Enforces admin rules, requires 1 approving review, dismisses stale reviews

### Terraform Output

```
$ terraform output
repository_name = "S25-core-course-labs"
repository_url = "https://github.com/escalopa/S25-core-course-labs"
repository_full_name = "escalopa/S25-core-course-labs"
default_branch = "master"
```

---

## Best Practices

The following Terraform best practices were applied in this project:

1. **Use a consistent file structure**: Each infrastructure component (Docker, Yandex Cloud, GitHub) has its own directory with `main.tf`, `variables.tf`, and `outputs.tf`.

2. **Never commit secrets**: Sensitive values (tokens, credentials) are passed via environment variables (`TF_VAR_*`) or `.tfvars` files excluded from version control via `.gitignore`.

3. **Use variables for configurable values**: All configurable parameters (container name, VM specs, repo settings) are defined as input variables with sensible defaults.

4. **Define outputs**: Key resource attributes are exported as outputs for easy reference and integration with other tools.

5. **Use `.gitignore` for Terraform**: State files (`*.tfstate`), variable files (`*.tfvars`), crash logs, and `.terraform/` directories are excluded from version control.

6. **Use `terraform import` for existing resources**: Instead of recreating resources, existing ones (like the GitHub repo) were imported into state using `terraform import`.

7. **Pin provider versions**: Provider versions are pinned using `~>` constraints (e.g., `~> 6.0`) to avoid unexpected breaking changes while still allowing minor updates.

8. **Mark sensitive variables**: Variables containing secrets (like tokens) are marked with `sensitive = true` to prevent them from being displayed in logs or plan output.

9. **Use `required_version`**: The minimum Terraform version is specified to ensure compatibility across environments.

10. **Keep state files secure**: State files may contain sensitive data and should never be committed. For team use, consider remote state backends (e.g., S3, Terraform Cloud).
