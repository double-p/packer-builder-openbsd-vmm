# packer-builder-openbsd-vmm
[Packer](https://packer.io/) builder plugin for OpenBSD's VMM

## Talk
Find my BSDCan 2019 slides in https://github.com/double-p/presentations/tree/master/BSDCan/2019

## jumpstart
```
make install
packer build examples/openbsd.json
packer build -var-file=examples/config.json examples/openbsd.json
```
More details in BUILD.md

## Example template

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
