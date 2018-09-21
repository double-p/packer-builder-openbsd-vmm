package main

import (
	"log"
	"github.com/hashicorp/packer/packer/plugin"
	vmm "./builder/openbsd-vmm"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		log.Fatal(err)
	}

	server.RegisterBuilder(new(vmm.Builder))
	server.Serve()
}
