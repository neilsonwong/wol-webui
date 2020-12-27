package main

// DeviceType enum
const (
	COMPUTER string = "COMPUTER"
	PHONE    string = "PHONE"
	OTHER    string = "OTHER"
)

// Device represents a
type Device struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Name string `json:"name"`
	IP   string `json:"ip"`
	Mac  string `json:"mac"`
}
