package openbsdvmm

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepOutDir struct {
	outputPath string
	cleanup    bool
	force      bool
}

func (step *stepOutDir) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	// Check if the output directory exists.
	if _, err := os.Stat(step.outputPath); err == nil {
		// If the build isn't forced, error out here.
		if !step.force {
			state.Put("error", fmt.Errorf("output directory already exists: %s", step.outputPath))
			return multistep.ActionHalt
		}

		// Forced build, so remove the directory.
		ui := state.Get("ui").(packer.Ui)
		ui.Say("Deleting previous output directory...")
		os.RemoveAll(step.outputPath)
	}

	// Mark that a cleanup is definately needed.
	step.cleanup = true

	// Create the output directory.
	if err := os.MkdirAll(step.outputPath, 0755); err != nil {
		state.Put("error", fmt.Errorf("output %s", step.outputPath))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (step *stepOutDir) Cleanup(state multistep.StateBag) {
	// Skip if no output directory was made in the first place.
	if !step.cleanup {
		return
	}

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Deleting output directory...")
	if err := os.RemoveAll(step.outputPath); err != nil {
		ui.Error("Unable to delete output directory: " + err.Error())
	}
}
