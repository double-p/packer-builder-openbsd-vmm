package openbsdvmm

import (
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

func CommHost() func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		sshHost := state.Get("ssh_host").(string)
		return sshHost, nil
	}
}

func SSHPort() func(multistep.StateBag) (int, error) {
	return func(state multistep.StateBag) (int, error) {
		config := state.Get("config").(*Config)
		return config.CommConfig.SSHPort, nil
	}
}
