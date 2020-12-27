package main

import (
	"fmt"
)

// Init queries the arp list then adds the missing entries into the device list
func Init() {
	fmt.Println("Init Devices from ARP Table")
	t := LoadARPTable()
	for ip, mac := range t {
		AddDevice(ip, mac)
	}
}

// ListDevices lists all devices
func ListDevices() []Device {
	return LoadDevices()
}

// GetDevice returns a device
func GetDevice(ID string) *Device {
	return GetDeviceByID(ID)
}

// Wakes the device by sending a magic packet
func WakeDevice(device Device) {
	Wake(device.Mac)
}
