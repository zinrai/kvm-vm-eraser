# kvm-vm-eraser

`kvm-vm-eraser` is a command-line tool designed to safely erase KVM (Kernel-based Virtual Machine) virtual machines and their associated disk images.

## Features

- Checks if the specified VM is running before attempting to erase it
- Undefines the VM using `virsh undefine`
- Deletes the associated disk image
- Requires sudo privileges for safety and to ensure proper cleanup
- Handles VMs with multiple attached disks, attempting to delete all of them
- Continues the erasure process even if some disk deletions fail, providing warnings for any failed deletions

## Notes

- This tool requires sudo privileges to run. It will prompt for your password if necessary.
- The tool will check if the VM is running before attempting to erase it. If the VM is running, it will display an error message and exit.
- Be extremely cautious when using this tool. It will permanently delete the VM and its associated disk image.
- Always ensure you have backups of important data before erasing a VM.
- The tool assumes that the VM name is unique and that there's only one disk associated with the VM.

## Installation

Build the tool:

```
$ go build
```

## Usage

Basic usage:

```
$ kvm-vm-eraser VM_NAME
```

Example:

```
$ kvm-vm-eraser my-virtual-machine
```

## License

This project is licensed under the [MIT License](./LICENSE).
