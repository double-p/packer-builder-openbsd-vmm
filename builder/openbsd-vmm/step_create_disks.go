package openbsdvmm

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepCreateDisks struct {
	outputPath string
	image      string
	format     string
	size       string
}

func (step *stepCreateDisks) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	var usedoas bool = false
	ui := state.Get("ui").(packer.Ui)
	path := filepath.Join(step.outputPath, step.image)

	command := []string{
		"create",
		//step.format + ":" + path,
		path,
		"-s",
		step.size,
	}
	ui.Say("Creating disk images...")
	if err := driver.VmctlCmd(usedoas, command...); err != nil {
		err := fmt.Errorf("Error creating disk image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("disk_image", path)

	return multistep.ActionContinue
}

func (step *stepCreateDisks) Cleanup(state multistep.StateBag) {
}
