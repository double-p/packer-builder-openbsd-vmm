#
source "openbsd-vmm" "alpine-base" {
    vm_name                        = "alpine-base-{{ isotime \"2006-01-02\" }}"
    vm_template                    = var.vm_template
    disk_format                    = var.disk_format
    disk_size                      = "10G"
    boot_device                    = "cdrom"
    cdrom                          = "/home/_vmd/_iso/alpine-virt-3.11.5-x86_64.iso"
    boot_wait                      = var.boot_wait
    boot_command                   = [
        "<enter><wait10>",
	"root<enter><wait>",
        "ifconfig eth0 up && udhcpc -i eth0<enter><wait10>",
        "wget http://{{ .HTTPIP }}:{{ .HTTPPort }}/alpine.autoinstall<enter><wait>",
        "wget http://{{ .HTTPIP }}:{{ .HTTPPort }}/authorized_keys<enter><wait>",
        "setup-alpine -ef alpine.autoinstall<enter><wait5>",
        "<wait30>y<enter>",
        "<wait120>",
	"mount /dev/vda2 /mnt<enter>",
	"mkdir -m 700 /mnt/root/.ssh && cp authorized_keys /mnt/root/.ssh/<enter>",
	"umount /dev/vda2<enter>",
        "/sbin/poweroff<enter>"
    ]

    gen_files_pattern              = "alpine"

    communicator                   = var.communicator
    ssh_agent_auth                 = var.ssh_agent_auth
    ssh_username                   = var.ssh_username

    shutdown_command               = "/sbin/poweroff"

    http_directory                 = var.http_directory
    log_directory                  = var.log_directory
    output_directory               = var.output_directory
}

build {
    sources                        = [ "source.openbsd-vmm.alpine-base" ]
}
