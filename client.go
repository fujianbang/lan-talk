package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	var tcpPort = 9000
	localIP := getLocalIP()
	log.Printf("[sys] get local ip: %s", localIP)

	localAddr := &net.UDPAddr{
		IP:   net.ParseIP(localIP),
		Port: 8888,
	}
	remoteAddr := &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 8888,
	}

	// udp broadcast
	udpConn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		log.Fatalf("[error] :%v\n", err)
		os.Exit(1)
	}
	defer udpConn.Close()


	// get tcp addr
	tcpAddr := &net.TCPAddr{
		IP:   net.ParseIP(localIP),
		Port: tcpPort,
	}
	_, err = udpConn.Write([]byte(tcpAddr.String()))
	if err != nil {
		log.Fatalf("[err] %v\n", err)
		os.Exit(1)
	}
	log.Printf("[sys] declare tcp addr: %s\n", tcpAddr)

	// dial tcp
	log.Println(tcpAddr.String())
	conn, err := waitConnectTcp(tcpAddr.String())

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("input message, type Q to quit")
		input, _ := inputReader.ReadString('\n')
		if trimInput := strings.TrimRight(input, "\n"); trimInput == "Q" {
			return
		}
		_, err = conn.Write([]byte(input))
	}
}

func getLocalIP() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("error: %v", err)
		return ""
	}

	for _, addr := range addrList {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

func waitConnectTcp(addr string) (conn net.Conn, err error) {
	for {
		time.Sleep(500 * time.Millisecond)
		log.Println("尝试建立tcp连接...")

		conn, err = net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("[error] %v", err)
		}
		return
	}
}
