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
	Msg struct {
		Header  Header
		Queries []Query
	}
)

func parseHeader(data []byte, pos int) (Header, int) {
	header := Header{}
	header.ID, pos = readUint16(data, pos)
	header.Flags, pos = readUint16(data, pos)
	header.QDCount, pos = readUint16(data, pos)
	header.ANCount, pos = readUint16(data, pos)
	header.NSCount, pos = readUint16(data, pos)
	header.ARCount, pos = readUint16(data, pos)

	// TODO(alex): remove.
	d, _ := json.Marshal(header)
	fmt.Printf("HEADER: %s\n", d)

	return header, pos
}

func parseQueries(header Header, data []byte, pos int) ([]Query, int) {
	queries := make([]Query, header.QDCount)
	for i := range queries {
		queries[i].QName, pos = readName(data, pos)
		queries[i].QType, pos = readQType(data, pos)
		queries[i].QClass, pos = readQClass(data, pos)

		// TODO(alex): remove.
		d, _ := json.Marshal(queries[i])
		fmt.Printf("QUERY %d: %s\n", i, d)
	}
	return queries, pos
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

	// TODO(alex): remove
	fmt.Printf("ADDR: %s\n", addr)
	fmt.Printf("RAW DATA: %b\n", data[0:n])

	parseMessage(data)
}
func parseMessage(data []byte) Msg {
	var pos int
	msg := Msg{}
	msg.Header, pos = parseHeader(data, pos)
	msg.Queries, pos = parseQueries(msg.Header, data, pos)
	return msg
}
