package openbsdvmm

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/hashicorp/packer/common"
	//"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
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
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"boot_command",
			},
		},
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
	if b.config.BootImage == "" {
		b.config.BootImage = "/bsd.rd"
	}
	if b.config.ImageName == "" {
		b.config.ImageName = "image-" + b.config.PackerBuildName
	}
	if b.config.OutDir == "" {
		b.config.OutDir = fmt.Sprintf("output-%s", b.config.PackerBuildName)
	}
	if b.config.DiskSize == "" {
		b.config.DiskSize = "5G"
	}
	switch b.config.DiskFormat {
	case "raw", "qcow2":
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
	b.config.bootWait, err = time.ParseDuration(b.config.RawBootWait)

	errs = packer.MultiErrorAppend(errs, isoErrs...)
	if len(errs.Errors) > 0 {
		return nil, errors.New(errs.Error())
	}

	return warnings, nil
}

// direct the workflow of creating the resulting artficat into "steppers"
func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	driver, err := b.newDriver()
	if err != nil {
		return nil, fmt.Errorf("Failed creating VMM driver: %s", err)
	}

	steps := []multistep.Step{}

	steps = append(steps, &stepOutDir{
		outputPath: b.config.OutDir,
		force:      b.config.PackerForce,
	})

	steps = append(steps, &stepCreateDisks{
		outputPath: b.config.OutDir,
		image:      b.config.ImageName,
		format:     b.config.DiskFormat,
		size:       b.config.DiskSize,
	})

	steps = append(steps, &stepLaunchVM{
		outputPath: b.config.OutDir,
		image:      b.config.ImageName,
		name:       b.config.VMName,
		mem:        b.config.RAMSize,
		kernel:     b.config.BootImage,
	})

	steps = append(steps, &common.StepHTTPServer{
		HTTPDir:     b.config.HTTPDir,
		HTTPPortMin: b.config.HTTPPortMin,
		HTTPPortMax: b.config.HTTPPortMax,
	})

	steps = append(steps, &stepBootCmd{
		cmd: b.config.FlatBootCommand(),
		ctx: b.config.ctx,
	})

	state := new(multistep.BasicStateBag)
	state.Put("driver", driver)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Run; step-wise if -debug/-on-error=ask
	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	if b.config.PackerDebug {
		b.runner = common.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	}
	b.runner.Run(context.Background(), state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	artifact := &VmmArtifact{
		imageDir: 	b.config.OutDir,
		imageName:	[]string{b.config.ImageName},
	}
	return artifact, nil
}

func (b *Builder) newDriver() (Driver, error) {
	// XXX: check doas.conf basics/existance
	doasbin := "/usr/bin/doas"
	// XXX: check VMD capable (see vagrant-openbsd-driver)
	vmctlbin := "/usr/sbin/vmctl"
	log := filepath.Join(b.config.OutDir+"/../", b.config.VMName+".log")
	driver := &vmmDriver{
		doas:    doasbin,
		logfile: log,
		vmctl:   vmctlbin,
	}
	return driver, nil
}
