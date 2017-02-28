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
	Query struct {
		QName  string
		QType  QType
		QClass QClass
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

func readQType(data []byte, pos int) (QType, int) {
	n, pos := readUint16(data, pos)
	return QType(n), pos
}

func readQClass(data []byte, pos int) (QClass, int) {
	n, pos := readUint16(data, pos)
	return QClass(n), pos
}

func readUint16(data []byte, ini int) (uint16, int) {
	fin := ini + 2
	return binary.BigEndian.Uint16(data[ini:fin]), fin
}

func readName(data []byte, pos int) (string, int) {
	var name string
	for data[pos] != 0 {
		labelLength := int(data[pos])
		pos++
		label := data[pos : pos+labelLength]
		pos += labelLength
		name += string(label) + "."
	}
	pos++
	return name, pos
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

	q := Query{}
	q.QName, pos = readName(data, pos)
	q.QType, pos = readQType(data, pos)
	q.QClass, pos = readQClass(data, pos)

	d, _ = json.Marshal(q)
	fmt.Printf("MSG: %s\n", d)
}
