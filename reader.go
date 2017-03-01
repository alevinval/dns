package dns

import (
	"bytes"
	"encoding/binary"
)

type Reader struct {
	data []byte
	i    int

	nameBuffer bytes.Buffer
}

func NewReader(b []byte) *Reader {
	return &Reader{data: b}
}

func (r *Reader) ReadMessage() *Msg {
	msg := &Msg{Header: Header{}}
	r.readHeader(&msg.Header)

	msg.Queries = make([]Query, msg.Header.QDCount)
	r.readQueries(msg.Queries)

	return msg
}

func (r *Reader) readQType() QType {
	return QType(r.readOctetPair())
}

func (r *Reader) readQClass() QClass {
	return QClass(r.readOctetPair())
}

func (r *Reader) readHeader(header *Header) {
	header.ID = r.readOctetPair()

	flags := r.readOctetPair()
	maskQR := uint16(1 << 15)
	maskOpCode := uint16(1 << 11)
	maskAA := uint16(1 << 10)
	maskTC := uint16(1 << 9)
	maskRD := uint16(1 << 8)
	maskRA := uint16(1 << 7)
	maskZ := uint16(7 << 4)
	maskRCode := uint16(16)

	header.QR = (flags & maskQR >> 15) == 1
	header.OpCode = OpCode(flags & maskOpCode >> 11)
	header.AA = (flags & maskAA >> 10) == 1
	header.TC = (flags & maskTC >> 9) == 1
	header.RD = (flags & maskRD >> 8) == 1
	header.RA = (flags & maskRA >> 7) == 1
	header.Z = byte(flags & maskZ)
	header.RCode = byte(flags & maskRCode)

	header.QDCount = r.readOctetPair()
	header.ANCount = r.readOctetPair()
	header.NSCount = r.readOctetPair()
	header.ARCount = r.readOctetPair()
}

func (r *Reader) readQueries(queries []Query) {
	for i := range queries {
		queries[i].QName = r.readName()
		queries[i].QType = r.readQType()
		queries[i].QClass = r.readQClass()
	}
}

func (r *Reader) readName() string {
	r.nameBuffer.Reset()
	for {
		currentByte := r.data[r.i]
		r.i++

		// On null termination, the label is finished.
		if currentByte == 0 {
			break
		}

		// Otherwise, current byte is the length of the label. Read it.
		ini := r.i
		r.i += int(currentByte)
		r.nameBuffer.Write(r.data[ini:r.i])
		r.nameBuffer.Write([]byte("."))
	}
	return r.nameBuffer.String()
}

func (r *Reader) readOctetPair() uint16 {
	ini := r.i
	r.i += 2
	return binary.BigEndian.Uint16(r.data[ini:r.i])
}
