package main

import (
	"flag"
	"log"
	"net"

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

	data := make([]byte, 512)
	for {
		n, peer, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("unexpected error reading from UDP: %s\n", err)
			continue
		}

		queryMsg, nUnpack, err := dns.UnpackMsg(data[:n], 0)
		if err != nil {
			log.Panicf("error unpacking message: %s\n", err)
		} else if nUnpack != n {
			log.Panicf("unpacked less bytes than what was read (read %d bytes, unpacked %d bytes)", n, nUnpack)
		}

		debug.PrintMessage(queryMsg)

		responseMsg := debug.FakeResponseMsg(queryMsg)
		response := dns.PackMsg(responseMsg)
		conn.WriteToUDP(response, peer)
	}
}
