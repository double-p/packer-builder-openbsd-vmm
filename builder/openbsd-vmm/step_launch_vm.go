package openbsdvmm

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepLaunchVM struct {
	name   string
	mem    string
	kernel string
	iso    string
}

func (step *stepLaunchVM) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	diskImage := state.Get("disk_image").(string)

	command := []string{
		"start",
		step.name,
		"-c",
		"-L",
		"-B",
		"net",
		"-i",
		"1",
		"-m",
		step.mem,
		"-b",
		step.kernel,
		"-d",
		diskImage}

	if step.iso != "" {
		command = append(command,
			"-r",
			step.iso,
		)
	}

	ui.Say("Bringing up VM...")
	if err := driver.Start(command...); err != nil {
		err := fmt.Errorf("Error bringing VM up: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (step *stepLaunchVM) Cleanup(state multistep.StateBag) {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	if err := driver.Stop(step.name); err != nil {
		e := fmt.Errorf("stopping vm (%s): %v", step.name, err)
		state.Put("error", e)
		ui.Error(e.Error())
	}
}
