#
source "openbsd-vmm" "centos-base" {
    vm_name                        = "centos-base"
    vm_template                    = "generic"
    disk_format                    = "qcow2"
    disk_size                      = "20G"
    cdrom                          = "/home/_vmd/_iso/CentOS-Stream-8-x86_64-20210907-dvd1.iso"
    boot_wait                      = "5s"
    boot_command                   = [
        "<esc><wait2>",
        "vmlinuz",
	" initrd=initrd.img",
	" inst.text",
	" nomodeset",
	" console=ttyS0,115200",
	" no_timer_check",
	" net.ifnames=0",
	" modprobe.blacklist=intel_pmc_core",
	" ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/centos.autoinstall",
	" ipv6.disable=1",
	"<enter>"
    ]

    gen_files_pattern              = "centos"

    communicator                   = "centos"
    ssh_agent_auth                 = "ssh"
    ssh_timeout                    = "1h"
    ssh_username                   = packer

    shutdown_command               = "sudo /sbin/halt -p"

    http_directory                 = "./_http"
    log_directory                  = "${var.home}/.log/packer"
    output_directory               = "/home/_vmd"
}

build {
    sources                        = [ "source.openbsd-vmm.centos-base" ]
}
