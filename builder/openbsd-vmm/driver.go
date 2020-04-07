package openbsdvmm

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/packer/common/bootcommand"
	"github.com/hashicorp/packer/packer"
)

type Driver interface {
	bootcommand.BCDriver
	VmctlCmd(args ...string) error
	Start(...string) error
	Stop(string) error
	GetTapIPAddress(string) (string, error)
	GetVMId(string) string
}

type vmmDriver struct {
	vmctl   string
	logfile string
	tty     io.WriteCloser
	console int
	ui      packer.Ui
}

func (d *vmmDriver) GetVMId(name string) string {
	var stdout bytes.Buffer
	cmd := exec.Command("vmctl", "status", name)
	cmd.Stdout = &stdout
	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("vmctl status error")
	}
	stdoutString := strings.TrimSpace(stdout.String())
	vmctl := regexp.MustCompile(`(\d+)`)
	resultarr := vmctl.FindAllStringSubmatch(stdoutString, -1)
	if resultarr == nil {
		return "VMAWOL"
	}
	return resultarr[0][1]
}

func (d *vmmDriver) VmctlCmd(args ...string) error {
	var stdout, stderr bytes.Buffer
	var cmd *exec.Cmd
	log.Printf("Executing command: vmctl %s", strings.Join(args, " "))
	cmd = exec.Command(d.vmctl, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("vmctl error")
	}
	log.Printf("stdout: %s", stdoutString)
	log.Printf("stderr: %s", stderrString)
	return err
}

// Start the VM and create a pipe to insert commands into the VM. (from packer-builder-vmm)
func (d *vmmDriver) Start(args ...string) error {
	logFile, err := os.OpenFile(d.logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	log.Printf("Executing vmctl: vmctl %s", strings.Join(args, " "))

	cmd := exec.Command(d.vmctl, args...)
	cmd.Env = append(os.Environ(),
		"TERM=vt220",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Create an stdin pipe that is used to issue commands.
	if d.tty, err = cmd.StdinPipe(); err != nil {
		return err
	}

	// Write the console output to the log file.
	go func() {
		defer stdout.Close()
		defer logFile.Close()

		_, _ = io.Copy(logFile, stdout)
	}()

	// Start up the VM.
	if err := cmd.Start(); err != nil {
		return err
	}

	// Give the VM a bit of time to start up.
	time.Sleep(3 * time.Second)
	return nil
}

func (d *vmmDriver) Stop(name string) error {
	cmd := exec.Command(d.vmctl, "stop", name)
	//err := cmd.Run()
	cmd.Run()
	return nil
}

func (d *vmmDriver) GetTapIPAddress(id string) (string, error) {
	var stdout bytes.Buffer
	vmId, _ := strconv.Atoi(id)
	vmName := fmt.Sprintf("vm%d", vmId)
	log.Printf("VM name: %s", vmName)

	// grab all available interfaces from group "tap"
	cmd := exec.Command("ifconfig", "tap")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("ifconfig error")
	}

	// parse interface(s) description and IPv4 addr
	stdoutString := strings.TrimSpace(stdout.String())
	log.Printf("ifconfig: %s", stdoutString)
	// XXX works on OpenBSD 6.6, but ugly
	vmctl := regexp.MustCompile(`description:\s(\w+\d+).*\n.*\n.*\n.*\n.*inet (\d+\.\d+\.\d+\.\d+)`)
	resultarr := vmctl.FindAllStringSubmatch(stdoutString, -1)
	// in case of multiple tap interfaces, loop into the result in order
	// to find the one we started
	for _, line := range resultarr {
		// [1] is the vmName
		// [2] is the IP
		if line[1] == vmName {
			return line[2], err
		}
	}
	err = fmt.Errorf("couldn't parse interface description")
	return "", err
}

//// interface Seq requires the following, not using it so far
// SendKey sends a key press.
func (d *vmmDriver) SendKey(key rune, action bootcommand.KeyAction) error {
	data := []byte{byte(key)}

	if len(data) != 0 {
		//log.Printf("Sending key '%q'", data)
		time.Sleep(100 * time.Millisecond)
		if _, err := d.tty.Write(data); err != nil {
			return err
		}
	}

	return nil
}

// SendSpecial sends a special character.
func (d *vmmDriver) SendSpecial(special string, action bootcommand.KeyAction) error {
	var data []byte

	switch special {
	case "bs":
		data = []byte("\b")
	case "enter":
		data = []byte("\n")
	case "esc":
		data = []byte("\x1b")
	case "tab":
		data = []byte("\t")
	}

	if len(data) != 0 {
		//log.Printf("Sending key '%q'", data)
		time.Sleep(100 * time.Millisecond)
		if _, err := d.tty.Write(data); err != nil {
			return err
		}
	}

	return nil
}

func (driver *vmmDriver) Flush() error {
	return nil
}
