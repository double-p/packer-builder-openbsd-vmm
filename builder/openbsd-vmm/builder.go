package hcloud

import (
	"fmt"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	"github.com/pkg/errors"
)

const BuilderID = "packer.openbsd-vmm"

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
	}, raws...)

	if err != nil {
		return nil, err
	}

	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, b.config.Comm.Prepare(&b.config.ctx)...)

	if b.config.ImageName == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("Missing image name"))
	}

	if len(errs.Errors) > 0 {
		return nil, errors.New(errs.Error())
	}

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	//client := openbsd-vmm.NewClient(openbsd-vmm.WithToken(b.config.Token))

	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("client", client)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&stepCreateSSHKey{},
		new(stepCreateServer),
		new(stepWaitForServer),
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      commHost,
			SSHConfig: sshConfig,
		},
		new(common.StepProvision),
		new(stepShutdown),
		new(stepPowerOff),
		new(stepCaptureImage),
		new(stepWaitForImage),
	}

	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(state)

	if rawErr, ok := state.GetOk("error"); ok {
		ui.Error(fmt.Sprintf("Got state error: %s", rawErr.(error)))
		return nil, rawErr.(error)
	}

	artifact := &Artifact{
		imageID:   state.Get("image_id").(int),
		imageName: state.Get("image_name").(string),
	}

	return artifact, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		b.runner.Cancel()
	}
}
