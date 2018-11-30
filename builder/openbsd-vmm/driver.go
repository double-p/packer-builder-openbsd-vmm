package openbsdvmm

import (
	"github.com/hashicorp/packer/packer"
)

type vmmDriver struct {
	doas string
	vmctl string
	logfile	string
	console int
	ui packer.Ui
}
