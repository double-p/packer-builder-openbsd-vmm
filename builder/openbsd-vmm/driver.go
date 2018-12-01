package openbsdvmm

import (
	"log"
	"os/exec"
	//"fmt"
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
	//var cmd *exec.Cmd
	if usedoas {
		log.Printf("Executing doas vmctl: %#v", args)
		args = append([]string{d.vmctl}, args...)
		cmd := exec.Command(d.doas, args...)
		err := cmd.Run()
		return err
	} else {
		log.Printf("Executing vmctl: %#v", args)
		cmd := exec.Command(d.vmctl, args...)
		err := cmd.Run()
		return err
	}
	//if _, ok := err.(*exec.ExitError); ok {
		//err = fmt.Errorf("vmctl error")
	//}
	return nil
}
