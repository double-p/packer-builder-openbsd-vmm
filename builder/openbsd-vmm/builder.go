package openbsdvmm

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
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
	if b.config.VMName == "" {
		b.config.VMName = "packer-" + b.config.PackerBuildName
	}
	if b.config.SourceImage == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("Missing source_image name"))
	}
	if b.config.ImageName == "" {
		b.config.ImageName = "image-" + b.config.PackerBuildName
	}
	if b.config.DiskSize == "" {
		b.config.DiskSize = "5G"
	}
	switch b.config.DiskFormat {
	case "raw,qcow2":
		// valid, use as is
	case "":
		b.config.DiskFormat = "raw"
	default:
		errs = packer.MultiErrorAppend(errs, errors.New("Unsupported disk_format name: "+b.config.DiskFormat))
	}
	if b.config.RAMSize == "" {
		b.config.RAMSize = "512M"
	}
	if b.config.Inet4 == "" {
		b.config.Inet4 = "dhcp" //XXX: some syntax check isIP4+prefix
	}
	if b.config.Inet4GW == "" {
		// "is ip4"?
	}
	if b.config.Inet6 == "" {
		b.config.Inet4 = "autoconf" //XXX: some syntax check isIP6+prefix
	}
	if b.config.Inet6GW == "" {
		// "is ip6"?
	}
	// XXX: DNS

	if len(errs.Errors) > 0 {
		return nil, errors.New(errs.Error())
	}

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	//client := openbsd-vmm.NewClient(openbsd-vmm.WithToken(b.config.Token))

	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	//state.Put("client", client)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		//&stepCreateSSHKey{},
		//new(stepCreateServer),
		//new(stepWaitForServer),
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			//Host:      commHost,
			//SSHConfig: sshConfig,
		},
		new(common.StepProvision),
		//new(stepShutdown),
		//new(stepPowerOff),
		//new(stepCaptureImage),
		//new(stepWaitForImage),
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
