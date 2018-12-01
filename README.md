# packer-builder-openbsd-vmm
[Packer](https://packer.io/) builder plugin for OpenBSD's VMM

# Building
```
go get github.com/hashicorp/packer/packer/plugin
go get github.com/pkg/errors
mkdir -p ~/.packer.d/plugins/
go build -o ~/.packer.d/plugins/packer-builder-openbsd-vmm
packer build examples/openbsd.json
packer build -var-file=examples/config.json examples/openbsd.json
```
(OpenBSD isnt on 1.11 yet, so no go.mod)

# Example template

```
{
  "builders": [
      {
          "type": "openbsd-vmm",
          "image_name": "some-image",
          "user_data": "",
          "ssh_username": "root"
      }
  ]
}
```

# Remarks
This is heavily based on https://github.com/m110/packer-builder-hcloud
