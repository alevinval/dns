package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 53})
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	data := make([]byte, 576)
	n, addr, _ := conn.ReadFromUDP(data)
	fmt.Printf("%s\n", addr)
	fmt.Printf("%b\n", data[0:n])
}
