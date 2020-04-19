#
source "openbsd-vmm" "centos-base" {
    vm_name                        = "centos-base-{{ isotime \"2006-01-02\" }}"
    vm_template                    = var.vm_template
    disk_format                    = var.disk_format
    disk_size                      = "20G"
    boot_device                    = "cdrom"
    cdrom                          = "/home/_vmd/_iso/CentOS-8.1.1911-x86_64-dvd1.iso"
    boot_wait                      = var.boot_wait
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
	"<enter>"
    ]

    gen_files_pattern              = "centos"

    communicator                   = var.communicator
    ssh_agent_auth                 = var.ssh_agent_auth
    ssh_timeout                    = "1h"
    ssh_username                   = var.ssh_username

    shutdown_command               = var.shutdown_command

    http_directory                 = var.http_directory
    log_directory                  = var.log_directory
    output_directory               = var.output_directory
}

build {
    sources                        = [ "source.openbsd-vmm.centos-base" ]
}
