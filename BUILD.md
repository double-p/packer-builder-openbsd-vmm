# Building plugin

Just a scratchboard for now..

## Packages
Needing a Golang eco system and some basics
````
pkg_add go-- packer-- git--
````

## Repo
````
git clone https://github.com/double-p/packer-builder-openbsd-vmm.git
````

## builds
Set ````GOPATH```` (default: ~/go), if the 1.4GB dependencies wont fit.

### OpenBSD < 6.6:
````
make
make install
````

### OpenBSD >= 6.6:
Apply the pull request #5 first.
https://github.com/double-p/packer-builder-openbsd-vmm/pull/5
Then
````
make
make install
````
Adding a version based switcher later this month.
