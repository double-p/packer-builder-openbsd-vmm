package openbsdvmm

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

const BuilderID = "packer.openbsd-vmm"

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	warnings, errs := b.config.Prepare(raws...)
	if errs != nil {
		return nil, warnings, errs
	}

	return nil, warnings, errs
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	driver, err := b.newDriver()
	if err != nil {
		return nil, fmt.Errorf("Failed creating VMM driver: %s", err)
	}

	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("driver", driver)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{}

	steps = append(steps,
		&stepOutDir{
			outputPath: b.config.OutDir,
			name:       b.config.VMName,
			format:     b.config.DiskFormat,
			force:      b.config.PackerForce,
		},
	)

	steps = append(steps,
		&stepCreateDisks{
			outputPath: b.config.OutDir,
			name:       b.config.VMName,
			format:     b.config.DiskFormat,
			baseImage:  b.config.DiskBase,
			size:       b.config.DiskSize,
		},
	)

	steps = append(steps,
		&stepStartVM{
			descr:    "OS installation",
			name:     b.config.VMName,
			mem:      b.config.MemorySize,
			bootdev:  b.config.BootDevice,
			kernel:   b.config.Boot,
			iso:      b.config.CdRom,
			template: b.config.VMTemplate,
		},
	)

	steps = append(steps,
		&common.StepHTTPServer{
			HTTPDir:     b.config.HTTPDir,
			HTTPPortMin: b.config.HTTPPortMin,
			HTTPPortMax: b.config.HTTPPortMax,
		},
	)

	steps = append(steps,
		&stepGenFiles{
			ctx: b.config.ctx,
		},
	)

	steps = append(steps,
		&stepBootCmd{
			BootWait: b.config.BootWait,
			cmd:      b.config.FlatBootCommand(),
			ctx:      b.config.ctx,
		},
	)

	// after install, boot from disk
	steps = append(steps,
		&stepStartVM{
			descr:    "provisioning",
			name:     b.config.VMName,
			mem:      b.config.MemorySize,
			iso:      b.config.CdRom,
			template: b.config.VMTemplate,
		},
	)

	steps = append(steps,
		&communicator.StepConnect{
			Config:    &b.config.CommConfig,
			Host:      CommHost(),
			SSHConfig: b.config.CommConfig.SSHConfigFunc(),
		},
	)

	steps = append(steps,
		&common.StepProvision{},
	)

	steps = append(steps,
		&stepShutdown{},
	)

	// Run; step-wise if -debug/-on-error=ask
	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	if b.config.PackerDebug {
		b.runner = common.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	}
	b.runner.Run(ctx, state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	// If we were interrupted or cancelled, then just exit.
	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("Build was cancelled.")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("Build was halted.")
	}

	artifact := &VmmArtifact{
		imageDir:  b.config.OutDir,
		imageName: []string{b.config.VMName},
	}
	return artifact, nil
}

func (b *Builder) newDriver() (Driver, error) {
	// XXX: check VMD capable (see vagrant-openbsd-driver)
	vmctlbin := "/usr/sbin/vmctl"
	log := filepath.Join(b.config.LogDir, b.config.VMName+".log")
	driver := &vmmDriver{
		logfile: log,
		vmctl:   vmctlbin,
	}
	return driver, nil
}
