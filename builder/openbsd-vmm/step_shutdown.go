package openbsdvmm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

// This step shuts down the machine. It first attempts to do so gracefully,
// but ultimately forcefully shuts it down if that fails.
//
// Uses:
//   communicator packer.Communicator
//   config *config
//   driver Driver
//   ui     packer.Ui
//   vm_id  string
//
// Produces:
//   <nothing>
type stepShutdown struct{}

func (step *stepShutdown) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	comm := state.Get("communicator").(packer.Communicator)
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	vmid := state.Get("vm_id").(string)

	if config.ShutdownCommand != "" {
		ui.Say("Gracefully halting virtual machine...")
		log.Printf("Executing shutdown command: %s", config.ShutdownCommand)
		cmd := &packer.RemoteCmd{Command: config.ShutdownCommand}
		if err := cmd.RunWithUi(ctx, comm, ui); err != nil {
			err := fmt.Errorf("Failed to send shutdown command: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

	} else {
		ui.Say("Halting the virtual machine...")
		if err := driver.Stop(vmid); err != nil {
			err := fmt.Errorf("Error stopping VM: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	// Wait for the machine to actually shut down
	log.Printf("Waiting max %s for shutdown to complete", config.ShutdownTimeout)
	for {
		halted := driver.GetVMId(vmid)

		if halted == "VMAWOL" {
			break
		}

		select {
		case <-time.After(config.ShutdownTimeout):
			err := errors.New("Timeout while waiting for machine to shutdown.")
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		default:
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("VM shut down.")
	return multistep.ActionContinue
}

func (step *stepShutdown) Cleanup(state multistep.StateBag) {}
