package networking

import (
	"fmt"
	"net"

	gowol "github.com/sabhiram/go-wol/wol"
)

// WakeByIP will wake a node by IP
func WakeByIP(ipAddr string) error {
	// Get MAC by IP
	mac := getMacFromIP(ipAddr)

	// wake by MAC
	return wake(mac)
}

func getMacFromIP(ipAddr string) string {
	return ""
}

func wake(macAddr string) error {
	bcastAddr := "255.255.255.255:9"
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)

	if err != nil {
		return err
	}

	// Build the magic packet.
	mp, err := gowol.New(macAddr)
	if err != nil {
		return err
	}

	// Grab a stream of bytes to send.
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}

	// Grab a UDP connection to send our packet of bytes.
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
	fmt.Printf("... Broadcasting to: %s\n", bcastAddr)
	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Magic packet sent successfully to %s\n", macAddr)
	return nil
}