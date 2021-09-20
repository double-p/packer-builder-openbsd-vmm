#
source "openbsd-vmm" "openbsd-dev" {
    vm_name                        = "openbsd-dev"
    vm_template                    = "generic"
    memory                         = "2G"
    disk_base                      = "/home/_vmd/openbsd-base.qcow2"
    disk_format                    = "qcow2"
    boot_wait                      = "5s"
    boot_command                   = [
        "<enter>",
        "mount -a -t ffs",
        "<enter>",
        "rm /etc/ssh/*key* /etc/iked/private/local.key /etc/iked/local.pub /etc/isakmpd/private/local.key /etc/isakmpd/local.pub /etc/soii.key",
        "<enter>",
        "sed -i /boot/d /etc/boot.conf",
        "<enter>",
        "echo dev.local > /etc/myname",
        "<enter>",
        "halt -p",
        "<enter>"
    ]

    ssh_agent_auth                 = "ssh"
    ssh_agent_auth                 = true
    ssh_username                   = "packer"

    shutdown_command               = "doas /sbin/halt -p"

    http_directory                 = "./_http"
    log_directory                  = "${var.home}/.log/packer"
    output_directory               = "/home/_vmd"
}

build {
    sources                        = [ "source.openbsd-vmm.openbsd-dev" ]

    # configure IP/DNS for vmm-running VM
    provisioner "shell" {
        inline = [
            "sleep 30",
            "doas su root -c \"echo 'inet 100.64.0.100 255.192.0.0 NONE' > /etc/hostname.vio0\"",
            "doas su root -c \"echo 'nameserver 100.64.0.1' > /etc/resolv.conf\"",
            "doas su root -c \"echo 'lookup file bind' >> /etc/resolv.conf\"",
            "doas su root -c \"echo 'search local' >> /etc/resolv.conf\"",
            "doas su root -c \"echo '100.64.0.1' > /etc/mygate\""
        ]
    }
}
