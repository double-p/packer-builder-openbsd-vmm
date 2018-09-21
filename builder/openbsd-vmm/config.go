package openbsd-vmm

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	ImageName   string `mapstructure:"image_name"`
	SourceImage string `mapstructure:"source_image"`

	UserData   string `mapstructure:"user_data"`

	ctx interpolate.Context
}
