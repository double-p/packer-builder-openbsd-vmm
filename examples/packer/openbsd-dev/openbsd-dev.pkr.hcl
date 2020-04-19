#
source "openbsd-vmm" "openbsd-dev" {
    vm_name                        = "openbsd-dev"
    vm_template                    = var.vm_template
    memory                         = "2G"
    disk_base                      = "/home/_vmd/openbsd-base-2020-04-19.qcow2"
    disk_format                    = var.disk_format
    boot_wait                      = var.boot_wait
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

    ssh_agent_auth                 = var.ssh_agent_auth
    ssh_username                   = "packer"

    shutdown_command               = "doas /sbin/halt -p"

    http_directory                 = var.http_directory
    log_directory                  = var.log_directory
    output_directory               = var.output_directory
}

build {
    sources                        = [ "source.openbsd-vmm.openbsd-dev" ]
    provisioner "shell" {
        inline = [
            "env TRUSTED_PKG_PATH='http://openbsd.local/%c/packages/%a/all' doas pkg_add go--",
            "touch ~/.hushlogin",
            "doas su root -c \"echo 'inet 100.64.0.100 255.192.0.0 NONE' > /etc/hostname.vio0\"",
            "doas su root -c \"echo 'nameserver 100.64.0.1' > /etc/resolv.conf\"",
            "doas su root -c \"echo 'lookup file bind' >> /etc/resolv.conf\"",
            "doas su root -c \"echo 'search local my.domain' >> /etc/resolv.conf\"",
            "doas su root -c \"echo '100.64.0.1' > /etc/mygate\"",
            "doas rm /etc/resolv.conf.tail",
            "doas install -c -o root -g wheel -m 664 /dev/null /etc/motd"
        ]
    }
}
