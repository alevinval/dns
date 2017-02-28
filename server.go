package main

import (
	"log"
	"net"

	"github.com/go-rfc/dns/debug"
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
	for {
		n, _, _ := conn.ReadFromUDP(data)
		p := parse.New(data[:n])
		msg := p.ReadMessage()
		debug.PrintMessage(msg)
	}
}
