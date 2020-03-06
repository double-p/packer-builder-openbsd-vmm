package openbsdvmm

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/mitchellh/go-homedir"

	//"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/pkg/errors"
)

const BuilderID = "packer.openbsd-vmm"

const (
	_DISK_QCOW2 = "qcow2"
	_DISK_RAW   = "raw"
)

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	if err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"boot_command",
			},
		},
	}, raws...); err != nil {
		return nil, nil, fmt.Errorf("decoding config: %v", err)
	}

	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, b.config.Comm.Prepare(&b.config.ctx)...)
	errs = packer.MultiErrorAppend(errs, b.config.HTTPConfig.Prepare(&b.config.ctx)...)

	if b.config.VMName == "" {
		b.config.VMName = "packer-" + b.config.PackerBuildName
	}

	if b.config.Bios == "" {
		b.config.Bios = "/etc/firmware/vmm-bios"
	}

	if b.config.OutDir == "" {
		b.config.OutDir = fmt.Sprintf("output-%s", b.config.PackerBuildName)
	}

	// DiskSize can be omitted if you're starting from a base image
	// as it'll use the same size as the base image.

	switch b.config.DiskFormat {
	case _DISK_RAW, _DISK_QCOW2:
		// valid, use as is
	case "":
		b.config.DiskFormat = _DISK_RAW
	default:
		errs = packer.MultiErrorAppend(errs, errors.New("Unsupported disk_format name: "+b.config.DiskFormat))
	}

	if b.config.DiskBase != "" && b.config.DiskFormat != _DISK_QCOW2 {
		errs = packer.MultiErrorAppend(errs, errors.New(
			"Cannot specify a base image without using qcow2 disk format"))
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

	var err error
	b.config.bootWait, err = time.ParseDuration(b.config.RawBootWait)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing bootwait time duration: %v", err)
	}

	if len(errs.Errors) > 0 {
		return nil, nil, errors.New(errs.Error())
	}

	return nil, nil, nil
}

// direct the workflow of creating the resulting artficat into "steppers"
func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	driver, err := b.newDriver()
	if err != nil {
		return nil, fmt.Errorf("Failed creating VMM driver: %s", err)
	}

	BiosPath, err := homedir.Expand(b.config.Bios)
	if err != nil {
		return nil, fmt.Errorf("failed expanding BIOS/kernel image path: %v", err)
	}

	isoPath, err := homedir.Expand(b.config.IsoImage)
	if err != nil {
		return nil, fmt.Errorf("failed expanding iso image path: %v", err)
	}

	var basePath string
	if b.config.DiskBase != "" {
		basePath, err = homedir.Expand(b.config.DiskBase)
		if err != nil {
			return nil, fmt.Errorf("failed expanding disk base path: %v", err)
		}
	}

	steps := []multistep.Step{}

	steps = append(steps, &stepOutDir{
		outputPath: b.config.OutDir,
		name:       b.config.PackerBuildName,
		force:      b.config.PackerForce,
	})

	steps = append(steps, &stepCreateDisks{
		outputPath: b.config.OutDir,
		name:       b.config.PackerBuildName,
		format:     b.config.DiskFormat,
		baseImage:  basePath,
		size:       b.config.DiskSize,
	})

	steps = append(steps, &stepLaunchVM{
		name:   b.config.VMName,
		mem:    b.config.RAMSize,
		kernel: BiosPath,
		iso:    isoPath,
	})

	steps = append(steps, &stepVMparams{
		name:   b.config.VMName,
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
		imageDir:  b.config.OutDir,
		imageName: []string{b.config.PackerBuildName},
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
