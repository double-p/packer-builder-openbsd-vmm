# variable definitions for all builds
variable "boot_wait"               { default = "10s" }
variable "communicator"            { default = "ssh" }
variable "disk_format"             { default = "qcow2" }
variable "http_directory"          { default = "/home/packer_user/.config/packer/_http" }
variable "log_directory"           { default = "/home/packer_user/.log/packer" }
variable "output_directory"        { default = "/home/_vmd" }
variable "shutdown_command"        { default = "/sbin/halt -p" }
variable "ssh_agent_auth"          { default = "true" }
variable "ssh_username"            { default = "root" }
variable "vm_template"             { default = "generic" }
