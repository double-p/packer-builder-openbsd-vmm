package openbsdvmm

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
)

type stepVMparams struct {
	name       string
}

func (step *stepVMparams) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	//ui := state.Get("ui").(packer.Ui)

	vmid := driver.GetVMId(step.name)

	state.Put("vm_id", vmid)
	return multistep.ActionContinue
}

func (step *stepVMparams) Cleanup(state multistep.StateBag) {}
