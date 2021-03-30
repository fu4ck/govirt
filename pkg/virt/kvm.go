package virt

import (
	"fmt"
	"github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"govirt/pkg/model"
	"log"
)

func GetVirtConn() (*libvirt.Connect, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	return conn, nil
}

func CreateXMLFile(instance *model.Instance, disk string, network string, configDrivePath string) (string, error) {
	domcfg := &libvirtxml.Domain{
		Type: "kvm",
		Name: instance.HostName,
		UUID: instance.Id,
		VCPU: &libvirtxml.DomainVCPU{
			Value:     16,
			Placement: "static",
		},
		CurrentMemory: &libvirtxml.DomainCurrentMemory{
			Unit:  "GiB",
			Value: uint(180),
		},
		CPU: &libvirtxml.DomainCPU{Mode: "host-passthrough",
			Topology: &libvirtxml.DomainCPUTopology{
				Sockets: 2,
				Threads: 2,
				Cores:   4,
			},
		},
		Memory: &libvirtxml.DomainMemory{
			Unit:  "GiB",
			Value: uint(180),
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch:    "x86_64",
				Machine: "pc-i440fx-xenial",
				Type:    "hvm",
			},
			BootDevices: []libvirtxml.DomainBootDevice{
				{
					Dev: "hd",
				},
			},
		},
		Devices: &libvirtxml.DomainDeviceList{
			Emulator: "/usr/bin/kvm",
			MemBalloon: &libvirtxml.DomainMemBalloon{
				Model: "virtio",
				Alias: &libvirtxml.DomainAlias{
					Name: "balloon0",
				},
				Address: &libvirtxml.DomainAddress{
					PCI: &libvirtxml.DomainAddressPCI{
						Domain:   helpUint(0),
						Bus:      helpUint(0),
						Slot:     helpUint(2),
						Function: helpUint(0),
					},
				},
			},
			Consoles: []libvirtxml.DomainConsole{
				{
					Source: &libvirtxml.DomainChardevSource{
						Pty: &libvirtxml.DomainChardevSourcePty{
							Path: "/dev/pts/7",
						},
					},
					Target: &libvirtxml.DomainConsoleTarget{
						Type: "serial",
						Port: helpUint(0),
					},
					Alias: &libvirtxml.DomainAlias{
						Name: "serial0",
					},
				},
			},
			Serials: []libvirtxml.DomainSerial{
				{
					Source: &libvirtxml.DomainChardevSource{
						Pty: &libvirtxml.DomainChardevSourcePty{
							Path: "/dev/pts/7",
						},
					},
					Target: &libvirtxml.DomainSerialTarget{
						Type: "isa-serial",
						Port: helpUint(0),
					},
					Alias: &libvirtxml.DomainAlias{
						Name: "serial0",
					},
				},
			},
			Controllers: []libvirtxml.DomainController{
				{
					Index: helpUint(0),
					Alias: &libvirtxml.DomainAlias{
						Name: "pci.0",
					},
					Type:  "pci",
					Model: "pci-root",
				},
			},
			Channels: []libvirtxml.DomainChannel{
				{
					Source: &libvirtxml.DomainChardevSource{
						UNIX: &libvirtxml.DomainChardevSourceUNIX{
							Mode: "bind",
							Path: "/var/lib/libvirt/qemu/org.qemu.guest_agent.0." + instance.Id + ".sock",
						},
					},
					Target: &libvirtxml.DomainChannelTarget{
						VirtIO: &libvirtxml.DomainChannelTargetVirtIO{
							Name:  "org.qemu.guest_agent.0",
							State: "connected",
						},
					},
					Alias: &libvirtxml.DomainAlias{
						Name: "channel0",
					},
					Address: &libvirtxml.DomainAddress{
						VirtioSerial: &libvirtxml.DomainAddressVirtioSerial{
							Controller: helpUint(0),
							Bus:        helpUint(0),
							Port:       helpUint(1),
						},
					},
				},
			},
		},
		Features: &libvirtxml.DomainFeatureList{
			APIC: &libvirtxml.DomainFeatureAPIC{},
			ACPI: &libvirtxml.DomainFeature{},
		},
		OnPoweroff: "destroy",
		OnCrash:    "destroy",
		OnReboot:   "restart",
		Resource:   &libvirtxml.DomainResource{Partition: "/machine"},
		Clock:      &libvirtxml.DomainClock{Offset: "localtime"},
	}
	// index for device
	index := 5

	configDrive := libvirtxml.DomainDisk{
		Device: "cdrom",
		Source: &libvirtxml.DomainDiskSource{
			File: &libvirtxml.DomainDiskSourceFile{
				File: configDrivePath, //TODO virt cdrom
			},
		},
		Target: &libvirtxml.DomainDiskTarget{
			Dev: "hdd",
			Bus: "ide",
		},
		ReadOnly: &libvirtxml.DomainDiskReadOnly{},
		Address: &libvirtxml.DomainAddress{
			Drive: &libvirtxml.DomainAddressDrive{
				Controller: helpUint(0),
				Bus:        helpUint(1),
				Target:     helpUint(0),
				Unit:       helpUint(1),
			},
		},
	}

	domcfg.Devices.Disks = append(domcfg.Devices.Disks, configDrive)

	// Add local disk
	for _, mount := range instance.LocalMount {

		LocalDrive := libvirtxml.DomainDisk{
			Device: "disk",
			Driver: &libvirtxml.DomainDiskDriver{
				Name:  "qemu",
				Type:  "qcow2",
				Cache: "default",
			},
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: mount.SourcePoint,
				},
			},
			Target: &libvirtxml.DomainDiskTarget{
				Dev: disk,
				Bus: "virtio",
			},

			Address: &libvirtxml.DomainAddress{
				PCI: &libvirtxml.DomainAddressPCI{
					Domain:   helpUint(0),
					Bus:      helpUint(0),
					Slot:     helpUint(uint(index)),
					Function: helpUint(0),
				},
			},
		}
		domcfg.Devices.Disks = append(domcfg.Devices.Disks, LocalDrive)
		//diskIndex++
		index++

	}

	// Add a vnc device to vm
	vncDevice := libvirtxml.DomainGraphic{
		VNC: &libvirtxml.DomainGraphicVNC{
			Listen:   "0.0.0.0",
			AutoPort: "no",
			Port:     5901,
			Keymap:   "en-us",
		},
	}

	instance.VNCUrl = instance.HostServer + ":" + "5901"

	domcfg.Devices.Graphics = append(domcfg.Devices.Graphics, vncDevice)

	xml, err := domcfg.Marshal()
	if err != nil {
		return "", err
	}

	return xml, nil
}
func helpUint(x uint) *uint { return &x }

func CreateVirualMachine(disk, configpath, network string) error {
	instance := model.NewInstance()
	xml, err := CreateXMLFile(instance, disk, configpath, network)
	if err != nil {
		log.Println(fmt.Sprintf("instance: %s define xml error", instance.Id))
		return err
	}
	conn, err := GetVirtConn()
	if err != nil {
		log.Println(fmt.Sprintf("instance can not get connection"))
		return err
	}

	domain, err := conn.DomainDefineXML(xml)

	if err != nil {
		log.Println(fmt.Sprintf("instance: %s define vm error", instance.Id))
		return err
	}

	err = domain.Create()

	if err != nil {
		log.Println(fmt.Sprintf("instance: %s start vm error", instance.Id))
		return err
	}
	return nil
}
