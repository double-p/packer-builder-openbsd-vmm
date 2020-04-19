# packer-builder-openbsd-vmm
[Packer](https://packer.io/) builder plugin for OpenBSD's VMM

## Talk
Find my BSDCan 2019 slides in https://github.com/double-p/presentations/tree/master/BSDCan/2019
Video is available at https://www.youtube.com/watch?v=GNmeFi3C1Xg

## jumpstart
Only for OpenBSD >=6.6
```
make install
packer build examples/openbsd.json
```
More details in BUILD.md

## bugs
Still some mad regexp about how to the find the connected tap(4) interface

If you find something, please use ``make vmb'' and include the log.

# Remarks
This is heavily based on https://github.com/m110/packer-builder-hcloud and
https://github.com/prep/packer-builder-vmm
