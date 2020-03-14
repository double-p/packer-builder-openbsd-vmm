package openbsdvmm

import (
	"time"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/common/shutdowncommand"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	common.HTTPConfig      `mapstructure:",squash"`
	bootcommand.BootConfig `mapstructure:",squash"`
	shutdowncommand.ShutdownConfig `mapstructure:",squash"`
	Comm                   communicator.Config `mapstructure:",squash"`
	RawBootWait            string              `mapstructure:"boot_wait"`
	bootWait               time.Duration       ``

	VMName     string `mapstructure:"vm_name" required:"true"`
	VMTemplate string `mapstructure:"vm_template" required:"true"`  // vmctl -t
	Console    bool   `mapstructure:"console"`     // vmctl -c
	BootDevice string `mapstructure:"boot_device"` // vmctl -B
	Boot       string `mapstructure:"boot"`        // vmctl -b
	CdRom      string `mapstructure:"cdrom"`       // vmctl -r
	DiskFormat string `mapstructure:"disk_format"` // vmctl create
	DiskBase   string `mapstructure:"disk_base"`   // vmctl create -b
	DiskSize   string `mapstructure:"disk_size"`   // vmctl create -s
	MemorySize string `mapstructure:"memory"`      // vmctl -m
	// not everybody lives in autoconf/DHCP; populate for hostname.vi0
	Inet4   string `mapstructure:"inet4"`       // hostname.if 'inet'
	Inet4GW string `mapstructure:"inet4gw"`     // mygate 'inet'
	Inet6   string `mapstructure:"inet6"`       // hostname.if 'inet6'
	Inet6GW string `mapstructure:"inet6gw"`     // mygate 'inet6'
	DNS     string `mapstructure:"nameservers"` // resolv.conf, comma-separated

	LogDir   string `mapstructure:"log_directory"`
	OutDir   string `mapstructure:"output_directory"`
	UserData string `mapstructure:"user_data"`

	ctx interpolate.Context
}
