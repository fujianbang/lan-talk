package test

import (
	"log"
	"net"
	"testing"
)

func TestLocalIP(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				log.Println(ipNet.IP.String())
			}
		}
	}
}
