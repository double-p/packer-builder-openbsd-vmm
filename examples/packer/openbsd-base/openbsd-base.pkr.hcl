#
source "openbsd-vmm" "openbsd-base" {
    vm_name                        = "openbsd-base"
    vm_template                    = "generic"
    disk_format                    = "qcow2"
    disk_size                      = "20G"
    boot                           = "/var/www/htdocs/openbsd/snapshots/amd64/bsd.rd"
    boot_device                    = "net"
    boot_wait                      = "5s"
    boot_command                   = [
        "A<enter><wait5>",
        "http://{{ .HTTPIP }}:{{ .HTTPPort }}/openbsd.autoinstall<enter>",
	"I<enter>"
    ]

    gen_files_pattern              = "openbsd"

    communicator                   = "ssh"
    ssh_agent_auth                 = true
    ssh_username                   = packer

    shutdown_command               = "doas /sbin/halt -p"

    http_directory                 = "./_http"
    log_directory                  = "${var.home}/.log/packer"
    output_directory               = "/home/_vmd"
}

build {
    sources                        = [ "source.openbsd-vmm.openbsd-base" ]
    provisioner "shell" {
        inline = [
            "sleep 180",
            "doas su root -c \"echo 'boot -s' >> /etc/boot.conf\"",
            "doas sed -i /openbsd/d /etc/hosts"
        ]
    }
}
