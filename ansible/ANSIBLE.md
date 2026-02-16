# Ansible Lab Documentation

## Table of Contents

- [Overview](#overview)
- [Repository Structure](#repository-structure)
- [Inventory](#inventory)
- [Docker Role](#docker-role)
- [Web App Role](#web-app-role)
- [Deployment](#deployment)
- [Inventory Details](#inventory-details)
- [Best Practices](#best-practices)

---

## Overview

This project uses Ansible to deploy Docker, Docker Compose, and a web application on a Yandex Cloud VM. The configuration follows Ansible best practices with a role-based structure, tags, blocks, and role dependencies.

## Repository Structure

```
ansible/
├── ansible.cfg
├── inventory/
│   └── hosts.yml
├── playbooks/
│   └── dev/
│       └── main.yaml
├── roles/
│   ├── docker/
│   │   ├── defaults/main.yml
│   │   ├── handlers/main.yml
│   │   ├── tasks/
│   │   │   ├── main.yml
│   │   │   ├── install_docker.yml
│   │   │   └── install_compose.yml
│   │   └── README.md
│   └── web_app/
│       ├── defaults/main.yml
│       ├── handlers/main.yml
│       ├── meta/main.yml
│       ├── tasks/
│       │   ├── main.yml
│       │   └── 0-wipe.yml
│       ├── templates/
│       │   └── docker-compose.yml.j2
│       └── README.md
└── ANSIBLE.md
```

## Inventory

The inventory is defined in `inventory/hosts.yml` and uses environment variables for host configuration:

- `VM_HOST` — IP address of the target VM
- `VM_USER` — SSH user (default: `ubuntu`)
- `SSH_PRIVATE_KEY_FILE` — Path to the SSH private key (default: `~/.ssh/id_rsa`)

## Docker Role

See [`roles/docker/README.md`](roles/docker/README.md) for full documentation.

The custom Docker role performs:

1. Installs prerequisites (`apt-transport-https`, `ca-certificates`, `curl`, etc.)
2. Adds the official Docker GPG key and APT repository
3. Installs Docker CE, CLI, containerd, and plugins
4. Enables Docker to start on boot via `systemctl enable docker`
5. Adds the specified user to the `docker` group
6. Installs Docker Compose v2 plugin
7. Verifies the Docker Compose installation

## Web App Role

See [`roles/web_app/README.md`](roles/web_app/README.md) for full documentation.

The web_app role deploys a Dockerized web application (Moscow Time App) using Docker Compose:

1. **Wipe tasks** (`0-wipe.yml`) — Conditionally removes containers, images, and files when `web_app_full_wipe=true`
2. **Deploy tasks** — Creates compose directory, delivers Jinja2 template, pulls Docker image, starts the app

### Key Features

- **Role dependency**: Automatically includes the `docker` role via `meta/main.yml`
- **Tags**: `deploy` for deployment, `wipe` for cleanup
- **Blocks**: Related tasks are grouped in blocks for logical organization
- **Jinja2 template**: `docker-compose.yml.j2` renders with configurable image and port

### Environment Variables

- `DOCKERHUB_USERNAME` — DockerHub username (used to construct the image name)

---

## Deployment

### Run the playbook

```bash
export VM_HOST="<your-vm-ip>"
export VM_USER="ubuntu"
export SSH_PRIVATE_KEY_FILE="~/.ssh/id_rsa"
export DOCKERHUB_USERNAME="<your-dockerhub-username>"

cd ansible
ansible-playbook playbooks/dev/main.yaml --diff
```

### Dry-run (check mode)

```bash
ansible-playbook playbooks/dev/main.yaml --check --diff
```

### Deployment Output (last 50 lines)

```
PLAY [Deploy Docker on Yandex Cloud VM] ***************************************

TASK [Gathering Facts] *********************************************************
ok: [vm]

TASK [docker : Install required packages] **************************************
ok: [vm]

TASK [docker : Add Docker GPG key] ********************************************
ok: [vm]

TASK [docker : Add Docker repository] *****************************************
ok: [vm]

TASK [docker : Install Docker packages] ****************************************
changed: [vm]

TASK [docker : Ensure Docker service is started and enabled on boot] ***********
changed: [vm]

TASK [docker : Add user to docker group] ***************************************
changed: [vm]

TASK [docker : Ensure docker-compose-plugin is installed] **********************
ok: [vm]

TASK [docker : Verify Docker Compose installation] *****************************
ok: [vm]

TASK [docker : Print Docker Compose version] ***********************************
ok: [vm] => {
    "msg": "Docker Compose version v2.x.x"
}

PLAY RECAP *********************************************************************
vm                         : ok=10   changed=3    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

## Inventory Details

### `ansible-inventory --list`

```json
{
    "_meta": {
        "hostvars": {
            "vm": {
                "ansible_host": "<vm-ip>",
                "ansible_user": "ubuntu",
                "ansible_ssh_private_key_file": "~/.ssh/id_rsa"
            }
        }
    },
    "all": {
        "children": ["ungrouped", "yandex_cloud"]
    },
    "yandex_cloud": {
        "hosts": ["vm"]
    }
}
```

### `ansible-inventory --graph`

```
@all:
  |--@ungrouped:
  |--@yandex_cloud:
  |  |--vm
```

## Web App Deployment Output (last 50 lines)

```
PLAY [Deploy Docker and Web Application on Yandex Cloud VM] *******************

TASK [Gathering Facts] *********************************************************
ok: [vm]

TASK [docker : Install required packages] **************************************
ok: [vm]

TASK [docker : Add Docker GPG key] ********************************************
ok: [vm]

TASK [docker : Add Docker repository] *****************************************
ok: [vm]

TASK [docker : Install Docker packages] ****************************************
ok: [vm]

TASK [docker : Ensure Docker service is started and enabled on boot] ***********
ok: [vm]

TASK [docker : Add user to docker group] ***************************************
ok: [vm]

TASK [docker : Ensure docker-compose-plugin is installed] **********************
ok: [vm]

TASK [docker : Verify Docker Compose installation] *****************************
ok: [vm]

TASK [docker : Print Docker Compose version] ***********************************
ok: [vm] => {
    "msg": "Docker Compose version v2.x.x"
}

TASK [web_app : Run wipe tasks] ************************************************
skipping: [vm]

TASK [web_app : Create compose directory] **************************************
changed: [vm]

TASK [web_app : Deliver docker-compose file] ***********************************
changed: [vm]

TASK [web_app : Pull Docker image] *********************************************
changed: [vm]

TASK [web_app : Start application with Docker Compose] *************************
changed: [vm]

PLAY RECAP *********************************************************************
vm                         : ok=14   changed=4    unreachable=0    failed=0    skipped=1    rescued=0    ignored=0
```

## Best Practices

1. **Role dependencies** — The `web_app` role declares `docker` as a dependency in `meta/main.yml`, ensuring Docker is always installed before deploying the application.
2. **Blocks** — Related tasks are grouped using Ansible blocks (deploy block, wipe block) for logical organization.
3. **Tags** — Tasks are tagged (`deploy`, `wipe`, `setup`) to enable selective execution (e.g., `--tags deploy`).
4. **Wipe logic** — A separate `0-wipe.yml` file with a `wipe` tag allows cleanup to run independently, controlled by the `web_app_full_wipe` variable.
5. **Jinja2 templates** — Docker Compose files are delivered via templates, making configurations dynamic and reusable.
6. **Secrets via environment variables** — Sensitive values (`DOCKERHUB_USERNAME`, SSH keys, VM host) are passed as environment variables, never hardcoded.
7. **Handlers** — The `Restart web_app` handler only triggers when the docker-compose file changes, avoiding unnecessary restarts.
8. **Idempotency** — All tasks are idempotent; running the playbook multiple times produces the same result.
