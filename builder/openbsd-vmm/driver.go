package openbsdvmm

import (
	"log"
	"os/exec"
	"fmt"
//	"github.com/hashicorp/packer/packer"
)

type Driver interface {
	VmctlCmd(bool, ...string) error
}

type vmmDriver struct {
	doas string
	vmctl string
	logfile	string
	console int
}

func (d *vmmDriver) VmctlCmd(usedoas bool, args ...string) error {
	log.Printf("Executing vmctl: %#v", args)
	cmd := exec.Command(d.vmctl, args...)
	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("vmctl error")
	}
	return err
}
