package main

import (
	"flag"
	"log"
	"net"
	"os"
)

var port = flag.Int("p", 8888, "udp port")

func main() {
	flag.Parse()

	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: *port,
	}
	tcpAddrStr := listenUDP(addr)

	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpAddrStr)
	if err != nil {
		log.Fatalf("[sys] %v", err)
		os.Exit(1)
	}
	log.Printf("正在建立TCP连接...(%s)", tcpAddr.String())
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("[sys] %v", err)
		os.Exit(1)
	}
	log.Printf("成功建立TCP连接...(%s)", tcpAddr.String())

	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			log.Fatalf("[error] %v", err)
			continue
		}
		go tcpHandler(conn)
	}
}

func listenUDP(addr *net.UDPAddr) string {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("[err] %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Printf("[sys] listening udp: %s", addr.String())

	//var tcpAddrBytes []byte
	tcpAddrBytes := make([]byte, 1024)
	_, remoteAddr, err := conn.ReadFromUDP(tcpAddrBytes)
	if err != nil {
		log.Fatalf("[err] %v\n", err)
	}
	log.Printf("[sys] remote addr: %s\n", remoteAddr.String())
	log.Printf("[sys] get tcp addr: %s\n", string(tcpAddrBytes))
	return string(tcpAddrBytes)
}

func tcpHandler(conn net.Conn) {
	log.Printf("连接上客户端：%s", conn.RemoteAddr().String())
	_, err := conn.Write([]byte("connection"))
	if err != nil {
		log.Fatalf("[error] %v", err)
		return
	}

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("[error] %v", err)
			continue
		}
		log.Printf("收到客户端的消息:%s", string(buf[:n]))
	}
}
