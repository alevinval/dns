package main

import (
	"flag"
	"log"
	"net"

	"bytes"

	"github.com/go-rfc/dns"
	"github.com/go-rfc/dns/debug"
)

var (
	port = flag.Int("port", 53, "binding port.")
)

func main() {
	flag.Parse()

	addr := &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: *port}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Panicf("server: cannot connect: %s", err)
	}
	defer conn.Close()
	log.Printf("server listening on %s", addr)

	data := make([]byte, 576)
	r := dns.NewReader(bytes.NewReader(data))
	for {
		_, peer, _ := conn.ReadFromUDP(data)
		msg, err := r.Read()
		if err == nil {
			debug.PrintMessage(msg)
			conn.WriteToUDP(data, peer)
		} else {
			log.Panicf("Error reading messages: %s\n", err)
		}

	}
}
