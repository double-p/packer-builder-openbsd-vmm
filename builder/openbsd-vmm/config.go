package openbsdvmm

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	common.HTTPConfig      `mapstructure:",squash"`
	common.ISOConfig       `mapstructure:",squash"`
	bootcommand.BootConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	VMName      string `mapstructure:"vm_name"`
	Console     bool   `mapstructure:"console"` // attach a console (to debug)
	BootImage   string `mapstructure:"boot_image"` //installNN.fs
	BootWait    string `mapstructure:"boot_wait"` //installNN.fs
	ImageName   string `mapstructure:"image_name"`  // artifact from build
	DiskSize    string `mapstructure:"disk_size"` // as vmctl -s
	DiskFormat  string `mapstructure:"disk_format"` // as vmctl create
	RAMSize     string `mapstructure:"ram_size"` // as vmctl -m
        // not everybody lives in autoconf/DHCP; populate for hostname.vi0
	Inet4       string `mapstructure:"inet4"` // hostname.if 'inet'
	Inet4GW     string `mapstructure:"inet4gw"` // mygate 'inet'
	Inet6       string `mapstructure:"inet6"` // hostname.if 'inet6'
	Inet6GW     string `mapstructure:"inet6gw"` // mygate 'inet6'
        DNS         string `mapstructure:"nameservers"` // resolv.conf, comma-separated
 
	OutDir      string `mapstructure:"output_directory"`
	UserData    string `mapstructure:"user_data"`

	ctx interpolate.Context
}
