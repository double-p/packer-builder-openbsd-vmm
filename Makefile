all: real

test:
	@go build -o testbuild
real:
	@go build -o ~/.packer.d/plugins/packer-builder-openbsd-vmm
