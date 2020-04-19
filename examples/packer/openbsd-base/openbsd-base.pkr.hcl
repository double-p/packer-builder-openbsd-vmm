#
source "openbsd-vmm" "openbsd-base" {
    vm_name                        = "openbsd-base-{{ isotime \"2006-01-02\" }}"
    vm_template                    = var.vm_template
    disk_format                    = var.disk_format
    disk_size                      = "20G"
    boot                           = "/var/www/htdocs/openbsd/snapshots/amd64/bsd.rd"
    boot_device                    = "net"
    boot_wait                      = var.boot_wait
    boot_command                   = [
        "http://{{ .HTTPIP }}:{{ .HTTPPort }}/openbsd.autoinstall<enter>",
	"I<enter>"
    ]

    gen_files_pattern              = "openbsd"

    communicator                   = var.communicator
    ssh_agent_auth                 = var.ssh_agent_auth
    ssh_username                   = "packer"

    shutdown_command               = "doas /sbin/halt -p"

    http_directory                 = var.http_directory
    log_directory                  = var.log_directory
    output_directory               = var.output_directory
}

build {
    sources                        = [ "source.openbsd-vmm.openbsd-base" ]
    provisioner "shell"   { inline = [
        "sleep 180",
        "doas su root -c \"echo 'boot -s' >> /etc/boot.conf\""
    ]}
}
