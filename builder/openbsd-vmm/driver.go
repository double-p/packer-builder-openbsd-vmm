package openbsdvmm

import (
	"bytes"
	"log"
	"os/exec"
	"fmt"
	"strings"
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
        var stdout, stderr bytes.Buffer
	var cmd *exec.Cmd
	if usedoas {
		args = append([]string{d.vmctl}, args...)
		args = append([]string{"ktrace"}, args...)
		log.Printf("Executing doas: %#v", args)
		cmd = exec.Command(d.doas, args...)
	} else {
		log.Printf("Executing vmctl: %#v", args)
		cmd = exec.Command(d.vmctl, args...)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("vmctl error")
	}
	log.Printf("stdout: %s", stdoutString)
	log.Printf("stderr: %s", stderrString)
	return err
}
