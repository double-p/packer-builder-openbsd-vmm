package openbsdvmm

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepStartVM struct {
	descr    string
	name     string
	mem      string
	bootdev  string
	kernel   string
	iso      string
	template string
}

func (step *stepStartVM) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	diskImage := state.Get("disk_image").(string)

	command := []string{
		"start",
		"-c",
		"-L",
		"-d",
		diskImage,
		"-t",
		step.template}
	if step.mem != "" {
		command = append(command,
			"-m",
			step.mem,
		)
	}
	if step.bootdev != "" {
		command = append(command,
			"-B",
			step.bootdev,
		)
	}
	if step.kernel != "" {
		command = append(command,
			"-b",
			step.kernel,
		)
	}
	if step.iso != "" {
		command = append(command,
			"-r",
			step.iso,
		)
	}
	command = append(command, step.name)

	ui.Say(fmt.Sprintf("Starting VM %s for %s...", step.name, step.descr))
	if err := driver.Start(command...); err != nil {
		err := fmt.Errorf("Error starting VM %s: %s", step.name, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	vmid := driver.GetVMId(step.name)
	hostIP, err := driver.GetTapIPAddress(vmid)
	if err != nil {
		state.Put("error", fmt.Errorf("Error getting hostIP: %s", err))
		return multistep.ActionHalt
	}
	vmIP := net.ParseIP(hostIP).To4()
	vmIP[3]++
	log.Printf("VM ID: %s", vmid)
	log.Printf("Host IP: %s", hostIP)
	log.Printf("VM IP: %s", vmIP.String())
	state.Put("step_descr", step.descr)
	state.Put("vm_id", vmid)
	state.Put("host_ip", hostIP)
	state.Put("ssh_host", vmIP.String())

	return multistep.ActionContinue
}

func (step *stepStartVM) Cleanup(state multistep.StateBag) {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	if err := driver.Stop(step.name); err != nil {
		err := fmt.Errorf("Error stopping VM (%s): %v", step.name, err)
		state.Put("error", err)
		ui.Error(err.Error())
	}
}
