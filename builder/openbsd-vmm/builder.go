package openbsdvmm

import (
	//"fmt"
	//"os/exec"
	"path/filepath"

	"github.com/hashicorp/packer/common"
	//"github.com/hashicorp/packer/helper/communicator"
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
        errs = packer.MultiErrorAppend(errs, b.config.HTTPConfig.Prepare(&b.config.ctx)...)
        warnings, isoErrs := b.config.ISOConfig.Prepare(&b.config.ctx)

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
	case "raw","qcow2":
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

        errs = packer.MultiErrorAppend(errs, isoErrs...)
	if len(errs.Errors) > 0 {
		return nil, errors.New(errs.Error())
	}

	return warnings, nil
}

// direct the workflow of creating the resulting artficat into "steppers"
func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	artifact := new(VmmArtifact)
/*
 instanciate driver
 steps:
*/
	steps := []multistep.Step{}
/*
 multistep collector array
 iso handling config
 empty disk config
*/
	steps = append(steps, &stepOutDir{
		outputPath: b.config.OutDir,
		force:      b.config.PackerForce,
	})
/*
 init internal http (autoinstall)
 bring in VM definition
 bootcommand/autoinstall
*/
	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Run; step-wise if -debug/-on-error=ask
	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	if b.config.PackerDebug {
		b.runner = common.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state) 
	}
	b.runner.Run(state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}
/*
 cast Artifact (wat?)
 return artifact
*/
	artifact.imageName = b.config.ImageName //faking artifact step
	artifact.imageSize = 123456321 //faking artifact step
	return artifact, nil
}

func (b *Builder) newDriver(ui packer.Ui) (vmmDriver, error) {
	// XXX: check doas.conf basics/existance
	doasbin := "/usr/bin/doas"
	// XXX: check VMD capable (see vagrant-openbsd-driver)
	vmctlbin := "/usr/sbin/vmctl"
        log := filepath.Join(b.config.OutDir, b.config.VMName + ".log")
        return vmmDriver {
		doas: doasbin,
		logfile: log,
		vmctl: vmctlbin,
		ui: ui,
	}, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		b.runner.Cancel()
	}
}
