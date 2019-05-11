SRC_FILES != /usr/bin/find . -type f -name '*.go'
PKG != /usr/local/bin/go list ./...

all: build

build:
	@go build

install: build
	@mkdir -p ~/.packer.d/plugins/
	@go build -o ~/.packer.d/plugins/packer-builder-openbsd-vmm

fmt:
	@gofmt -e -l -s $(SRC_FILES) | grep "go" && echo "gofmt -s -d on above" || exit 0

vet:
	@go vet -all $(PKG)

test:
	@go test -v -timeout 60s $(PKG)

clean:
	@rm -f packer-builder-openbsd-vmm

uninstall: clean
	@rm -f ~/.packer.d/plugins/packer-builder-openbsd-vmm

