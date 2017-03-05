package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/go-rfc/dns"
)

var (
	dumpFile = flag.String("dump", "dump.bin", "dump file path.")
	port     = flag.Int("port", 53, "binding port.")
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

	f, err := os.OpenFile(*dumpFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Panicf("error opening dump file: %s", err)
	}
	defer f.Close()

	data := make([]byte, 576)
	unpacker := dns.NewUnpacker()
	for {
		n, peer, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("unexpected error reading from UDP: %s\n", err)
			continue
		}

		unpacker.Reset(data[:n])
		msg, nUnpack, err := unpacker.Unpack()
		if nUnpack != n {
			log.Panicf("unpacked less bytes than what was read (read %d bytes, unpacked %d bytes)", n, nUnpack)
		}
		if err != nil {
			log.Panicf("error unpacking message: %s\n", err)
		}

		w := dns.NewWriter(msg)
		w.WriteTo(f)

		log.Printf("served request for %q from %s\n", msg.Queries[0].QName, peer.String())
		conn.WriteToUDP(data, peer)
	}
}
