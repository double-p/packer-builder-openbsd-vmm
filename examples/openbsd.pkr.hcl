variable "log_directory"        { default = "/home/packer_user/.log/packer" }
variable "output_directory"     { default = "/home/packer_user/.local/share/packer" }
variable "http_directory"       { default = "/home/packer_user/.config/packer/autoinstall" }
variable "vm_template"          { default = "generic" }
variable "disk_format"          { default = "qcow2" }
variable "boot_device"          { default = "net" }
variable "boot"                 { default = "/var/www/htdocs/openbsd.tristero.se/snapshots/amd64/bsd.rd" }
variable "boot_wait"            { default = "10s" }
variable "ssh_username"         { default = "packer_user" }
variable "ssh_agent_auth"	{ default = "true" }
variable "trusted_pkg_path"     { default = "http://192.168.255.1/pub/OpenBSD/%c/packages/%a/all" }
variable "shutdown_command"     { default = "doas /sbin/halt -p" }

source "openbsd-vmm" "openbsd" {
    vm_name          = "openbsd"
    vm_template      = "${var.vm_template}"
    memory           = "2G"
    disk_size        = "20G"
    disk_format      = "${var.disk_format}"
    boot_device      = "${var.boot_device}"
    boot             = "${var.boot}"
    boot_wait        = "${var.boot_wait}"
    boot_command     = [
        "http://{{ .HTTPIP }}:{{ .HTTPPort }}/autoinstall<enter>",
	"I<enter>"
    ]
    log_directory    = "${var.log_directory}"
    output_directory = "${var.output_directory}"
    http_directory   = "${var.http_directory}"
    ssh_username     = "${var.ssh_username}"
    ssh_agent_auth   = "${var.ssh_agent_auth}"
    communicator     = "ssh"
    shutdown_command = "${var.shutdown_command}"
}

build {
    sources          = [ "source.openbsd-vmm.openbsd" ]
    #provisioner "breakpoint" { note = "Debug" }
    provisioner "shell" { inline = [
	"sleep 300",
	"env TRUSTED_PKG_PATH='${var.trusted_pkg_path}' doas pkg_add unzip"
    ]}
}
