# Ansible Lab Documentation

## Table of Contents

- [Overview](#overview)
- [Repository Structure](#repository-structure)
- [Inventory](#inventory)
- [Docker Role](#docker-role)
- [Deployment](#deployment)
- [Inventory Details](#inventory-details)

---

## Overview

This project uses Ansible to deploy Docker and Docker Compose on a Yandex Cloud VM. The configuration follows Ansible best practices with a role-based structure.

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
│   └── docker/
│       ├── defaults/main.yml
│       ├── handlers/main.yml
│       ├── tasks/
│       │   ├── main.yml
│       │   ├── install_docker.yml
│       │   └── install_compose.yml
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

## Deployment

### Run the playbook

```bash
export VM_HOST="<your-vm-ip>"
export VM_USER="ubuntu"
export SSH_PRIVATE_KEY_FILE="~/.ssh/id_rsa"

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
