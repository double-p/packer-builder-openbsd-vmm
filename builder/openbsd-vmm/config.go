//go:generate mapstructure-to-hcl2 -type Config

package openbsdvmm

import (
	"fmt"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/common/shutdowncommand"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

const (
	_DISK_QCOW2 = "qcow2"
	_DISK_RAW   = "raw"
)

type Config struct {
	common.PackerConfig            `mapstructure:",squash"`
	common.HTTPConfig              `mapstructure:",squash"`
	bootcommand.BootConfig         `mapstructure:",squash"`
	shutdowncommand.ShutdownConfig `mapstructure:",squash"`

	CommConfig communicator.Config `mapstructure:",squash"`

	VMName     string `mapstructure:"vm_name"      required:"true"`
	VMTemplate string `mapstructure:"vm_template"  required:"true"` // vmctl -t
	Console    bool   `mapstructure:"console"`                      // vmctl -c
	BootDevice string `mapstructure:"boot_device"`                  // vmctl -B
	Boot       string `mapstructure:"boot"`                         // vmctl -b
	CdRom      string `mapstructure:"cdrom"`                        // vmctl -r
	DiskFormat string `mapstructure:"disk_format"`                  // vmctl create
	DiskBase   string `mapstructure:"disk_base"`                    // vmctl create -b
	DiskSize   string `mapstructure:"disk_size"`                    // vmctl create -s
	MemorySize string `mapstructure:"memory"`                       // vmctl -m

	LogDir   string `mapstructure:"log_directory"`
	OutDir   string `mapstructure:"output_directory"`
	UserData string `mapstructure:"user_data"`

	ctx interpolate.Context
}

func (c *Config) Prepare(raws ...interface{}) ([]string, error) {

	err := config.Decode(c, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &c.ctx,
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
	errs = packer.MultiErrorAppend(errs, c.CommConfig.Prepare(&c.ctx)...)
	errs = packer.MultiErrorAppend(errs, c.HTTPConfig.Prepare(&c.ctx)...)
	errs = packer.MultiErrorAppend(errs, c.ShutdownConfig.Prepare(&c.ctx)...)
	errs = packer.MultiErrorAppend(errs, c.BootConfig.Prepare(&c.ctx)...)

	if c.VMName == "" {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("VM name must be specified (var: vm_name)"))
	}

	if c.VMTemplate == "" {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("VM template must be specified (var: vm_template)"))
	}

	if c.OutDir == "" {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("Output directory must be specified (var: output_directory)"))
	}

	switch c.DiskFormat {
	case _DISK_RAW, _DISK_QCOW2:
	// use default raw format if not specified
	case "":
		c.DiskFormat = _DISK_RAW
	default:
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("Unsupported disk_format name: "+c.DiskFormat+", must be either raw or qcow2"))
	}

	if c.DiskBase == "" && c.DiskSize == "" {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("Disk size must be specified when not using base image (var: disk_size)"))
	}

	if c.DiskBase != "" && c.DiskFormat != _DISK_QCOW2 {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("Cannot use "+c.DiskFormat+" with base image, only qcow2 format is supported"))
	}

	if c.CommConfig.Type == "ssh" {
		if c.CommConfig.User() == "" {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("SSH Username must be specified (var: ssh_username)"))
		}

		if !c.CommConfig.SSHAgentAuth && c.CommConfig.Password() == "" && c.CommConfig.SSHPrivateKeyFile == "" {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("SSH authentication method must be specified (vars: ssh_agent_auth, ssh_password, ssh_private_key_file)"))
		}

		if (c.CommConfig.SSHAgentAuth &&
			(c.CommConfig.SSHPassword != "" || c.CommConfig.SSHPrivateKeyFile != "")) ||
			(c.CommConfig.SSHPassword != "" && c.CommConfig.SSHPrivateKeyFile != "") {
			errs = packer.MultiErrorAppend(errs,
				fmt.Errorf("Only one SSH authentication method is supported (vars: ssh_agent_auth, ssh_password, ssh_private_key_file)"))

		}
	} else if c.CommConfig.Type != "none" {
		errs = packer.MultiErrorAppend(errs,
			fmt.Errorf("Only ssh or none communicator is supported (var: communicator)"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	return nil, nil
}
