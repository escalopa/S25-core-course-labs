# Web App Role

This role deploys a Dockerized web application using Docker Compose.

## Requirements

- Ansible 2.9+
- Ubuntu 20.04 / 22.04 / 24.04
- Docker installed (handled automatically via `docker` role dependency)

## Role Variables

| Variable                  | Default                                              | Description                                      |
|---------------------------|------------------------------------------------------|--------------------------------------------------|
| `web_app_image`           | `{{ lookup('env', 'DOCKERHUB_USERNAME') }}/moscow-time-app:latest` | Docker image to deploy |
| `web_app_port`            | `5000`                                               | Host port to expose the application              |
| `web_app_container_name`  | `moscow-time-app`                                    | Name of the Docker container                     |
| `web_app_compose_dir`     | `/opt/web_app`                                       | Directory for the docker-compose file            |
| `web_app_full_wipe`       | `false`                                              | Enable full wipe (remove container, image, files)|

## Dependencies

- `docker` role (automatically included via `meta/main.yml`)

## Tags

| Tag      | Description                              |
|----------|------------------------------------------|
| `deploy` | Deploy/update the web application        |
| `wipe`   | Remove the application and all its files |

## Example Playbook

```yaml
- hosts: all
  become: true
  roles:
    - role: web_app
```

### Deploy only

```bash
ansible-playbook playbooks/dev/main.yaml --tags deploy
```

### Wipe application

```bash
ansible-playbook playbooks/dev/main.yaml --tags wipe -e web_app_full_wipe=true
```
