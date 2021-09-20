package main

import (
	"log"

	vmm "github.com/double-p/packer-builder-openbsd-vmm/builder/openbsd-vmm"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		log.Fatal(err)
	}

	server.RegisterBuilder(new(vmm.Builder))
	server.Serve()
}
