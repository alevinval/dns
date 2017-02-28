package parse

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type MessageReader struct {
	data []byte
	pos  int
}

func New(b []byte) *MessageReader {
	return &MessageReader{data: b}
}

func (mr *MessageReader) Read() *Msg {
	msg := &Msg{}
	mr.readMessage()
	return msg
}

func (mr *MessageReader) readQType() QType {
	n := mr.readUint16()
	return QType(n)
}

func (mr *MessageReader) readQClass() QClass {
	n := mr.readUint16()
	return QClass(n)
}

func (mr *MessageReader) readUint16() uint16 {
	ini := mr.pos
	mr.pos += 2
	return binary.BigEndian.Uint16(mr.data[ini:mr.pos])
}

func (mr *MessageReader) readMessage() Msg {
	msg := Msg{}
	msg.Header = mr.readHeader()
	msg.Queries = mr.readQueries(msg.Header)
	return msg
}

func (mr *MessageReader) readHeader() Header {
	header := Header{}
	header.ID = mr.readUint16()
	header.Flags = mr.readUint16()
	header.QDCount = mr.readUint16()
	header.ANCount = mr.readUint16()
	header.NSCount = mr.readUint16()
	header.ARCount = mr.readUint16()

	// TODO(alex): remove.
	d, _ := json.Marshal(header)
	fmt.Printf("HEADER: %s\n", d)

	return header
}

func (mr *MessageReader) readQueries(header Header) []Query {
	queries := make([]Query, header.QDCount)
	for i := range queries {
		queries[i].QName = mr.readName()
		queries[i].QType = mr.readQType()
		queries[i].QClass = mr.readQClass()

		// TODO(alex): remove.
		d, _ := json.Marshal(queries[i])
		fmt.Printf("QUERY %d: %s\n", i, d)
	}
	return queries
}

func (mr *MessageReader) readName() string {
	var name string
	for mr.data[mr.pos] != 0 {
		labelLength := int(mr.data[mr.pos])
		mr.pos++
		label := mr.data[mr.pos : mr.pos+labelLength]
		mr.pos += labelLength
		name += string(label) + "."
	}
	mr.pos++
	return name
}
