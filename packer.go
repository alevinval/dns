package dns

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Unpacker struct {
	data []byte
	pos  int

	msg        *Msg
	nameBuffer bytes.Buffer
}

func NewUnpacker(b []byte) *Unpacker {
	return &Unpacker{msg: &Msg{Header: Header{}}, data: b}
}

func (r *Unpacker) Unpack() (msg *Msg, n int, err error) {
	err = r.unpackHeader()
	if err != nil {
		return nil, 0, err
	}

	err = r.unpackQueries()
	if err != nil {
		return nil, 0, err
	}

	return r.msg, r.pos, io.EOF
}

func (r *Unpacker) Reset(b []byte) {
	r.pos = 0
	r.data = b
}

func (r *Unpacker) readQType() QType {
	return QType(r.unpackOctetPair())
}

func (r *Unpacker) readQClass() QClass {
	return QClass(r.unpackOctetPair())
}

func (r *Unpacker) unpackHeader() (err error) {
	if len(r.data) < 12 {
		return io.ErrShortBuffer
	}
	h := Header{}
	h.ID = r.unpackOctetPair()

	flags := r.unpackOctetPair()
	h.QR = (flags & maskQR >> maskQROffset) == 1
	h.OpCode = OpCode(flags & maskOpCode >> maskOpCodeOffset)
	h.AA = (flags & maskAA >> maskAAOffset) == 1
	h.TC = (flags & maskTC >> maskTCOffset) == 1
	h.RD = (flags & maskRD >> maskRDOffset) == 1
	h.RA = (flags & maskRA >> maskRAOffset) == 1
	h.Z = byte(flags & maskZ >> maskZOffset)
	h.RCode = byte(flags & maskRCode)

	h.QDCount = r.unpackOctetPair()
	h.ANCount = r.unpackOctetPair()
	h.NSCount = r.unpackOctetPair()
	h.ARCount = r.unpackOctetPair()

	r.msg.Header = h
	return nil
}

func (r *Unpacker) unpackQueries() (err error) {
	queries := make([]Query, r.msg.Header.QDCount)
	r.msg.Queries = queries
	for i := range queries {
		qName, err := r.readName()
		if err != nil {
			return err
		}
		if r.pos+2 >= len(r.data) {
			return io.ErrShortBuffer
		}
		queries[i].QName = qName
		queries[i].QType = r.readQType()
		queries[i].QClass = r.readQClass()
	}
	return nil
}

func (r *Unpacker) readName() (string, error) {
	r.nameBuffer.Reset()
	for {
		currentByte, err := r.unpackByte()
		if err != nil {
			return "", err
		}

		// On null termination, the label is finished.
		if currentByte == 0 {
			break
		}

		// Otherwise, current byte is the length of the label. Read it.
		ini := r.pos
		r.pos += int(currentByte)
		if r.pos >= len(r.data) {
			return "", io.ErrShortBuffer
		}
		r.nameBuffer.Write(r.data[ini:r.pos])
		r.nameBuffer.Write([]byte("."))
	}
	return r.nameBuffer.String(), nil
}

func (r *Unpacker) unpackOctetPair() uint16 {
	ini := r.pos
	r.pos += 2
	return binary.BigEndian.Uint16(r.data[ini:r.pos])
}

func (r *Unpacker) unpackByte() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, io.ErrShortBuffer
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}
