package main

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/go-rfc/dns"
	"github.com/go-rfc/dns/debug"
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
		n, peer, _ := conn.ReadFromUDP(data)
		r := dns.NewReader(data[:n])
		msg := r.ReadMessage()
		debug.PrintMessage(msg)

		w := dns.NewWriter(msg)
		w.Write()

		// TODO: remove this! create a dns-dump cli to dump packets into raw binary files.
		fmt.Printf("%b\n", data[:n])
		fmt.Printf("%b\n", w.Bytes())
		println(bytes.Equal(data[:n], w.Bytes()))

		conn.WriteToUDP(data, peer)
	}
}
