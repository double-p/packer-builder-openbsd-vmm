package openbsdvmm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type stepBootCmd struct {
	cmd string
	ctx interpolate.Context
}

type bootCommandTemplateData struct {
	HTTPIP   string
	HTTPPort uint
}

func (step *stepBootCmd) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	httpPort := state.Get("http_port").(uint)

        hostIp, err := driver.GetTapIPAddress(config.VMName)
        ui.Say(fmt.Sprintf("%s with Host HTTPD on %s:%d", config.VMName, hostIp, httpPort))
        common.SetHTTPIP(hostIp)
        step.ctx.Data = &bootCommandTemplateData{
                hostIp,
                httpPort,
        }

	ui.Say(fmt.Sprintf("boot_wait is (%s).", config.bootWait.String()))
	if int64(config.bootWait) > 0 {
		time.Sleep(time.Duration(config.bootWait))
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
	ui.Say("Waiting 15s for bootcommand to finish...") //XXX debug
	time.Sleep(15 * time.Second)
	return multistep.ActionContinue
}

func (step *stepBootCmd) Cleanup(state multistep.StateBag) {
}
