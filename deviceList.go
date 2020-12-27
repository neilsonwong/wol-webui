package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/uuid"
)

const deviceFile = "devices.json"

// LoadDevices loads the list of devices from file
func LoadDevices() []Device {
	file, _ := ioutil.ReadFile(deviceFile)
	var devices []Device
	_ = json.Unmarshal([]byte(file), &devices)
	return devices
}

// AddDevice adds an entry into the list of devices if it does not already exist
func AddDevice(ip string, mac string) {
	fmt.Printf("ip: %s, mac: %s\n", ip, mac)
	devices := LoadDevices()

	if i := findDeviceByIP(devices, ip); i != -1 {
		devices[i].Mac = mac
	} else {
		guid := uuid.New()
		devices = append(devices, Device{
			ID:   guid.String(),
			Name: "unknown",
			Kind: OTHER,
			IP:   ip,
			Mac:  mac,
		})
	}

	updateDeviceList(devices)
}

// AddNewDevice adds a new device that may not be found in the arp
func AddNewDevice(device Device) {
	devices := LoadDevices()

	if findDeviceByIP(devices, device.IP) != -1 {
		log.Printf("Attempt to readd device with id: %s, ip: %s, mac: %s", device.ID, device.IP, device.Mac)
		return
	}

	// add the new entry
	guid := uuid.New()
	device.ID = guid.String()
	devices = append(devices, device)
	updateDeviceList(devices)
}

func UpdateDevice(device Device) {
	devices := LoadDevices()

	if i := findDeviceByID(devices, device.ID); i != -1 {
		devices[i] = device
	} else {
		log.Printf("Attempt to update nonexistant device with id: %s", device.ID)
	}

	updateDeviceList(devices)
}

func updateDeviceList(device []Device) {
	file, _ := json.MarshalIndent(device, "", " ")
	_ = ioutil.WriteFile(deviceFile, file, 0644)
}

// GetDeviceByID gets the device by id
func GetDeviceByID(ID string) *Device {
	devices := LoadDevices()
	if i := findDeviceByID(devices, ID); i != -1 {
		return &devices[i]
	}
	return nil
}

func findDeviceByIP(devices []Device, IP string) int {
	for i := 0; i < len(devices); i++ {
		if devices[i].IP == IP {
			return i
		}
	}
	return -1
}

func findDeviceByID(devices []Device, ID string) int {
	for i := 0; i < len(devices); i++ {
		if devices[i].ID == ID {
			return i
		}
	}
	return -1
}

func init() {

}
