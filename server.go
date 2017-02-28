package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type (
	Header struct {
		ID      uint16
		Flags   uint16
		QDCount uint16
		ANCount uint16
		NSCount uint16
		ARCount uint16
	}
)

func readHeader(data []byte, pos int) (Header, int) {
	header := Header{}
	header.ID, pos = readUint16(data, pos)
	header.Flags, pos = readUint16(data, pos)
	header.QDCount, pos = readUint16(data, pos)
	header.ANCount, pos = readUint16(data, pos)
	header.NSCount, pos = readUint16(data, pos)
	header.ARCount, pos = readUint16(data, pos)
	return header, pos
}

func readUint16(data []byte, ini int) (uint16, int) {
	fin := ini + 2
	return binary.BigEndian.Uint16(data[ini:fin]), fin
}

func main() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 53})
	if err != nil {
		log.Panicf("server: cannot connect: %s", err)
	}
	defer conn.Close()

	data := make([]byte, 576)
	n, addr, _ := conn.ReadFromUDP(data)
	fmt.Printf("ADDR: %s\n", addr)
	fmt.Printf("RAW DATA: %b\n", data[0:n])

	var pos int
	header, pos := readHeader(data, pos)

	d, _ := json.Marshal(header)
	fmt.Printf("HEADER: %s\n", d)
}
