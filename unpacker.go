package dns

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	headerLen    = 12
	octetPairLen = 2
)

type Unpacker struct {
	buffer []byte
	msg    *Msg
	i      int

	nameBuffer bytes.Buffer
}

func NewUnpacker() *Unpacker {
	return &Unpacker{}
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

	return r.msg, r.i, nil
}

func (r *Unpacker) Reset(b []byte) {
	r.buffer = b
	r.msg = &Msg{Header: Header{}}
	r.i = 0
}

func (r *Unpacker) readQType() QType {
	return QType(r.unpackOctetPair())
}

func (r *Unpacker) readQClass() QClass {
	return QClass(r.unpackOctetPair())
}

func (r *Unpacker) unpackHeader() (err error) {
	if len(r.buffer) < headerLen {
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
		if r.i+octetPairLen >= len(r.buffer) {
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
		ini := r.i
		r.i += int(currentByte)
		if r.i >= len(r.buffer) {
			return "", io.ErrShortBuffer
		}
		r.nameBuffer.Write(r.buffer[ini:r.i])
		r.nameBuffer.Write([]byte("."))
	}
	return r.nameBuffer.String(), nil
}

func (r *Unpacker) unpackOctetPair() uint16 {
	ini := r.i
	r.i += octetPairLen
	return binary.BigEndian.Uint16(r.buffer[ini:r.i])
}

func (r *Unpacker) unpackByte() (byte, error) {
	if r.i >= len(r.buffer) {
		return 0, io.ErrShortBuffer
	}
	b := r.buffer[r.i]
	r.i++
	return b, nil
}
