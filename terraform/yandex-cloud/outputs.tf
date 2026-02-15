output "vm_id" {
  description = "ID of the created VM"
  value       = yandex_compute_instance.vm.id
}

output "vm_name" {
  description = "Name of the created VM"
  value       = yandex_compute_instance.vm.name
}

output "vm_external_ip" {
  description = "External IP address of the VM"
  value       = yandex_compute_instance.vm.network_interface[0].nat_ip_address
}

output "vm_internal_ip" {
  description = "Internal IP address of the VM"
  value       = yandex_compute_instance.vm.network_interface[0].ip_address
}

output "subnet_id" {
  description = "ID of the created subnet"
  value       = yandex_vpc_subnet.subnet.id
}
