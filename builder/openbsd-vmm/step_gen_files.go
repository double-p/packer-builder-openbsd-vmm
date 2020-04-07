package openbsdvmm

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/template/interpolate"
)

type stepGenFiles struct {
	ctx interpolate.Context
}

type genFilesTemplateData struct {
	VMName   string
	HTTPIP   string
	HTTPPort int
}

func scanLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func (step *stepGenFiles) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	httpPort := state.Get("http_port").(int)
	hostIP := state.Get("host_ip").(string)
	VMName := config.VMName

	step.ctx.Data = &genFilesTemplateData{
		VMName,
		hostIP,
		httpPort,
	}

	log.Printf("Generating files...")

	substVars := func(path string, fileinfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !!fileinfo.IsDir() {
			return nil
		}

		matched, err := filepath.Match(VMName + "*.pkr.in", fileinfo.Name())

		if matched {
			lines, err := scanLines(path)
			if err != nil {
				state.Put("error", fmt.Errorf("Error reading input file: %s", err))
				return err
			}

			newfile, err := os.OpenFile(strings.TrimSuffix(path, ".pkr.in"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				state.Put("error", fmt.Errorf("Error writing output file: %s", err))
				return err
			}

			defer newfile.Close()

			writer := bufio.NewWriter(newfile)

			for _, line := range lines {
				newline, err := interpolate.Render(line, &step.ctx)
				writer.WriteString(newline + "\n")
				if err != nil {
					state.Put("error", fmt.Errorf("Error rendering line: %s", err))
					return err
				}
			}
			return writer.Flush()
		}
		return err
	}

	if err := filepath.Walk(config.HTTPDir, substVars); err != nil {
		state.Put("error", fmt.Errorf("Error generating files: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (step *stepGenFiles) Cleanup(state multistep.StateBag) {}
