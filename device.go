package main

// DeviceType enum
const (
	COMPUTER string = "COMPUTER"
	PHONE    string = "PHONE"
	OTHER    string = "OTHER"
)

// Device represents a
type Device struct {
	id   string
	kind string
	name string
	ip   string
	mac  string
}
