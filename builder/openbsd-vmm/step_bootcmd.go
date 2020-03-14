package openbsdvmm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type stepBootCmd struct {
	cmd      string
	BootWait time.Duration
	ctx      interpolate.Context
}

type bootCommandTemplateData struct {
	HTTPIP   string
	HTTPPort int
}

func (step *stepBootCmd) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	step_descr := state.Get("step_descr").(string)
	httpPort := state.Get("http_port").(int)
	hostIP := state.Get("host_ip").(string)

	log.Printf("HTTP IP/port: %s:%d", hostIP, httpPort)
	common.SetHTTPIP(hostIP)
	step.ctx.Data = &bootCommandTemplateData{
		hostIP,
		httpPort,
	}

	if int64(step.BootWait) > 0 {
		ui.Say(fmt.Sprintf("Waiting %s before %s starts...", step.BootWait.String(), step_descr))
		time.Sleep(time.Duration(step.BootWait))
		select {
		case <-time.After(step.BootWait):
			break
		case <-ctx.Done():
			return multistep.ActionHalt
		}
	}
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
	ui.Say(fmt.Sprintf("Waiting until %s is completed and VM %s shutdown...", step_descr, config.VMName))
	for {
		halted := driver.GetVMId(config.VMName)
		if halted == "VMAWOL" {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return multistep.ActionContinue
}

func (step *stepBootCmd) Cleanup(state multistep.StateBag) {}
