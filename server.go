package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-rfc/dns/parse"
)

func main() {
	addr := &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 53}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Panicf("server: cannot connect: %s", err)
	}
	defer conn.Close()
	log.Printf("server listening on %s", addr)

	data := make([]byte, 576)
	n, addr, _ := conn.ReadFromUDP(data)

	// TODO(alex): remove
	fmt.Printf("ADDR: %s\n", addr)
	fmt.Printf("RAW DATA: %b\n", data[0:n])

	p := parse.New(data)
	p.ReadMessage()
}
