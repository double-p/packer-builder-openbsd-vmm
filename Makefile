SRC_FILES=$(shell find . -type f -name '*.go')
PKG=$(shell go list ./... )

all: build

build:
	@go build

install: build
	@mkdir -p ~/.packer.d/plugins/
	@go build -o ~/.packer.d/plugins/packer-builder-openbsd-vmm

fmt:
	@gofmt -e -l -s $(SRC_FILES) |grep ".*" && exit 2 || exit 0

vet:
	@go vet -all $(PKG)

test:
	@go test -v -timeout 60s $(PKG)

clean:
	@rm -f packer-builder-openbsd-vmm

uninstall: clean
	@rm -f ~/.packer.d/plugins/packer-builder-openbsd-vmm

