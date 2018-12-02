package openbsdvmm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/template/interpolate"
)

type stepBootCmd struct {
	wait   time.Duration
	cmd    string
	ctx    interpolate.Context
}

func (step *stepBootCmd) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	if step.wait > 0 {
		ui.Say(fmt.Sprintf("Waiting %s for boot...", step.wait.String()))
		select {
		case <-time.After(step.wait):
			break
		case <-ctx.Done():
			return multistep.ActionHalt
		}
	}

	ui.Say("Typing the boot command...")
	command, err := interpolate.Render(step.cmd, &step.ctx)
	if err != nil {
		state.Put("error", fmt.Errorf("Error preparing boot command: %s", err))
		return multistep.ActionHalt
	}

	seq, err := bootcommand.GenerateExpressionSequence(command)
	if err != nil {
		state.Put("error", fmt.Errorf("Error generating boot command: %s", err))
		return multistep.ActionHalt
	}

	if err := seq.Do(ctx, driver); err != nil {
		state.Put("error", fmt.Errorf("Error running boot command: %s", err))
		return multistep.ActionHalt
	}
	return multistep.ActionContinue
}

func (step *stepBootCmd) Cleanup(state multistep.StateBag) {
}
