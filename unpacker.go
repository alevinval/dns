package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	MaxNameLen  = 255
	MaxLabelLen = 63

	headerLen    = 12
	octetPairLen = 2
)

var (
	ErrLabelTooLong = errors.New("a label must be 63 octets or less")
	ErrNameTooLong  = errors.New("a name must be 255 octets or less")
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
		label, n, err := unpackLabel(r.buffer, r.i)
		r.i += n
		if err == io.EOF {
			return r.nameBuffer.String(), nil
		} else if err != nil {
			return "", err
		}
		r.nameBuffer.WriteString(label + ".")
		if r.nameBuffer.Len() > MaxNameLen {
			return "", ErrNameTooLong
		}
	}
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

func unpackLabel(b []byte, offset int) (label string, n int, err error) {
	if !checkBounds(b, offset, offset+1) {
		return "", 0, io.ErrShortBuffer
	}

	// Current byte indicates the length of the label.
	// If its a null byte, label is over.
	currentByte := b[offset]
	if currentByte == 0 {
		return "", 1, io.EOF
	}
	n++
	offset++

	// Check if its a pointer.
	isPointer := currentByte>>6 == 3
	if isPointer {

		// Compute the offset to the pointer.
		offset = int((currentByte&64)<<8 + b[offset])
		currentByte = b[offset]
		offset++
	}

	// Check if the label has valid length.
	endOffset := offset + int(currentByte)
	if endOffset-offset > MaxLabelLen {
		return "", 0, ErrLabelTooLong
	}
	if !checkBounds(b, offset, endOffset) {
		return "", 0, io.ErrShortBuffer
	}

	// Return the label.
	if !isPointer {
		n += endOffset - offset
	}
	return string(b[offset:endOffset]), n, nil
}

// Check if begin and end are within bounds of a byte slice.
func checkBounds(b []byte, begin, end int) bool {
	return len(b) >= begin && len(b[begin:]) >= end-begin
}
