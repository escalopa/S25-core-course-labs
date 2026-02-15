# Docker Role

This role installs and configures Docker and Docker Compose on Ubuntu systems.

## Requirements

- Ansible 2.9+
- Ubuntu 20.04 / 22.04 / 24.04
- Root or sudo access on the target machine

## Role Variables

| Variable                 | Default                  | Description                                      |
|--------------------------|--------------------------|--------------------------------------------------|
| `docker_user`            | `{{ ansible_user }}`     | User to add to the `docker` group                |
| `docker_compose_version` | `latest`                 | Docker Compose version (installed via apt plugin) |
| `docker_edition`         | `ce`                     | Docker edition (`ce` for Community Edition)       |
| `docker_packages`        | See `defaults/main.yml`  | List of Docker packages to install                |

## Tasks

1. **install_docker.yml** — Adds Docker GPG key and repository, installs Docker CE, enables the service on boot, adds the user to the `docker` group.
2. **install_compose.yml** — Installs `docker-compose-plugin` and verifies the installation.

## Handlers

- **Restart docker** — Restarts the Docker daemon and reloads systemd.

## Example Playbook

```yaml
- hosts: all
  become: true
  roles:
    - role: docker
```

## Testing

Run a dry-run to verify changes before applying:

```bash
ansible-playbook playbooks/dev/main.yaml --check --diff
```
