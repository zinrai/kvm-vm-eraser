package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Domain struct {
	Devices Devices `xml:"devices"`
}

type Devices struct {
	Disks []Disk `xml:"disk"`
}

type Disk struct {
	Type   string `xml:"type,attr"`
	Device string `xml:"device,attr"`
	Source Source `xml:"source"`
}

type Source struct {
	File string `xml:"file,attr"`
}

func main() {
	if !isVirshAvailable() {
		fmt.Println("Error: 'virsh' command not found. Please ensure libvirt-clients is installed.")
		os.Exit(1)
	}

	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Usage: kvm-vm-eraser VM_NAME")
		os.Exit(1)
	}

	vmName := flag.Arg(0)

	if isVMRunning(vmName) {
		fmt.Printf("Error: The VM '%s' is currently running. Please stop it before attempting to erase.\n", vmName)
		fmt.Println("You can stop the VM using the command: sudo virsh shutdown", vmName)
		os.Exit(1)
	}

	if err := eraseVM(vmName); err != nil {
		fmt.Printf("Failed to erase VM: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully erased VM %s and deleted its disk(s).\n", vmName)
}

func isVirshAvailable() bool {
	_, err := exec.LookPath("virsh")
	return err == nil
}

func isVMRunning(vmName string) bool {
	cmd := exec.Command("sudo", "virsh", "list", "--name", "--state-running")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to get list of running VMs: %v\n", err)
		return true
	}

	runningVMs := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, vm := range runningVMs {
		if vm == vmName {
			return true
		}
	}
	return false
}

func eraseVM(vmName string) error {
	diskPaths, err := getVMImagePaths(vmName)
	if err != nil {
		return fmt.Errorf("failed to get VM image paths: %v", err)
	}

	undefineCmd := exec.Command("sudo", "virsh", "undefine", vmName)
	if err := undefineCmd.Run(); err != nil {
		return fmt.Errorf("failed to undefine VM: %v", err)
	}

	for _, diskPath := range diskPaths {
		rmCmd := exec.Command("sudo", "rm", "-f", diskPath)
		if err := rmCmd.Run(); err != nil {
			fmt.Printf("Warning: failed to delete disk %s: %v\n", diskPath, err)
		} else {
			fmt.Printf("Deleted disk: %s\n", diskPath)
		}
	}

	return nil
}

func getVMImagePaths(vmName string) ([]string, error) {
	cmd := exec.Command("sudo", "virsh", "dumpxml", vmName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get VM XML: %v", err)
	}

	var domain Domain
	if err := xml.Unmarshal(output, &domain); err != nil {
		return nil, fmt.Errorf("failed to parse VM XML: %v", err)
	}

	var diskPaths []string
	for _, disk := range domain.Devices.Disks {
		if disk.Device == "disk" && disk.Type == "file" {
			diskPaths = append(diskPaths, disk.Source.File)
		}
	}

	if len(diskPaths) == 0 {
		return nil, fmt.Errorf("no disk images found for VM %s", vmName)
	}

	return diskPaths, nil
}
