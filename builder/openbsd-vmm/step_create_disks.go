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
	name       string
	format     string
	size       string
	baseImage  string
}

func (step *stepCreateDisks) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	var usedoas bool = false
	ui := state.Get("ui").(packer.Ui)
	path := filepath.Join(step.outputPath, step.name)

	// >= 6.6 format : vmctl [-v] create [-b base | -i disk] [-s size] disk
	command := []string{
		"create",
		"-s",
		step.size,
	}

	if step.baseImage != "" {
		command = append(command,
			"-b",
			step.baseImage)
	}

	command = append(command,
		step.format+":"+path)

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

func (step *stepCreateDisks) Cleanup(state multistep.StateBag) {}
