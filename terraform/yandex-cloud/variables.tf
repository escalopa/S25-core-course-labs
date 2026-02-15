variable "yc_token" {
  description = "Yandex Cloud OAuth token"
  type        = string
  sensitive   = true
}

variable "yc_cloud_id" {
  description = "Yandex Cloud ID"
  type        = string
}

variable "yc_folder_id" {
  description = "Yandex Cloud Folder ID"
  type        = string
}

variable "yc_zone" {
  description = "Yandex Cloud availability zone"
  type        = string
  default     = "ru-central1-a"
}

variable "vm_name" {
  description = "Name of the VM instance"
  type        = string
  default     = "lab4-vm"
}

variable "image_id" {
  description = "ID of the boot disk image (Ubuntu 22.04 LTS)"
  type        = string
  default     = "fd8autg36kchufhej85b" # Ubuntu 22.04 LTS
}

variable "network_id" {
  description = "ID of an existing Yandex VPC network"
  type        = string
  default     = "enpef4gkfef33crv2auj"
}

variable "vm_user" {
  description = "Username to create on the VM for SSH access"
  type        = string
  default     = "ubuntu"
}

variable "ssh_public_key_path" {
  description = "Path to the SSH public key file"
  type        = string
  default     = "~/.ssh/devops-course.pub"
}
