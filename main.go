package main

import (
	"fmt"
	"os"

	vmm "github.com/double-p/packer-builder-openbsd-vmm/builder/openbsd-vmm"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(vmm.Builder))
	pps.SetVersion(version.NewPluginVersion("0.8", "", ""))
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
