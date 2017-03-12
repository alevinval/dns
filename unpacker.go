package dns

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	MaxNameLen  = 255
	MaxLabelLen = 63

	lenHeader = 12
	lenUint16 = 2
)

var (
	ErrLabelTooLong = errors.New("a label must be 63 octets or less")
	ErrNameTooLong  = errors.New("a name must be 255 octets or less")
)

type Unpacker struct {
	buffer []byte
	msg    *Msg
	i      int
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
	qType, n, _ := unpackUint16(r.buffer, r.i)
	r.i += n
	return QType(qType)
}

func (r *Unpacker) readQClass() QClass {
	return QClass(r.unpackOctetPair())
}

func (r *Unpacker) unpackHeader() (err error) {
	h, n, err := unpackHeader(r.buffer, r.i)
	if err != nil {
		return err
	}
	r.msg.Header = *h
	r.i += n
	return nil
}

func (r *Unpacker) unpackQueries() (err error) {
	queries := make([]Query, r.msg.Header.QDCount)
	r.msg.Queries = queries
	for i := range queries {
		qName, n, err := unpackName(r.buffer, r.i)
		if err != nil {
			return err
		}
		r.i += n
		queries[i].QName = qName

		// TODO: pending refactor on readQType and readQClass
		if r.i+2*lenUint16 >= len(r.buffer) {
			return io.ErrShortBuffer
		}
		queries[i].QType = r.readQType()
		queries[i].QClass = r.readQClass()
	}
	return nil
}

func (r *Unpacker) unpackOctetPair() uint16 {
	ini := r.i
	r.i += lenUint16
	return binary.BigEndian.Uint16(r.buffer[ini:r.i])
}

func unpackHeader(b []byte, offset int) (h *Header, n int, err error) {
	if !checkBounds(b, offset, offset+lenHeader) {
		return nil, 0, io.ErrShortBuffer
	}
	iniOffset := offset
	h = &Header{}
	h.ID, n, err = unpackUint16(b, offset)
	offset += n
	if err != nil {
		return nil, 0, err
	}

	var flags uint16
	flags, n, err = unpackUint16(b, offset)
	offset += n
	if err != nil {
		return nil, 0, err
	}

	h.QR = (flags & maskQR >> maskQROffset) == 1
	h.OpCode = OpCode(flags & maskOpCode >> maskOpCodeOffset)
	h.AA = (flags & maskAA >> maskAAOffset) == 1
	h.TC = (flags & maskTC >> maskTCOffset) == 1
	h.RD = (flags & maskRD >> maskRDOffset) == 1
	h.RA = (flags & maskRA >> maskRAOffset) == 1
	h.Z = byte(flags & maskZ >> maskZOffset)
	h.RCode = byte(flags & maskRCode)

	h.QDCount, n, err = unpackUint16(b, offset)
	offset += n

	h.ANCount, n, err = unpackUint16(b, offset)
	offset += n
	if err != nil {
		return nil, 0, err
	}

	h.NSCount, n, err = unpackUint16(b, offset)
	offset += n
	if err != nil {
		return nil, 0, err
	}

	h.ARCount, n, err = unpackUint16(b, offset)
	offset += n
	if err != nil {
		return nil, 0, err
	}

	return h, offset - iniOffset, nil

}

func unpackUint16(b []byte, offset int) (r uint16, n int, err error) {
	end := offset + 1
	if !checkBounds(b, offset, end) {
		return 0, 0, io.ErrShortBuffer
	}
	return uint16(b[end]) | uint16(b[offset])<<8, 2, nil
}

func unpackName(b []byte, offset int) (name string, n int, err error) {
	var ln int
	var label string
	for {

		// Unpack a label and advance offset and read bytes.
		label, ln, err = unpackLabel(b, offset)
		offset += ln
		n += ln

		// Check for errors.
		if err == io.EOF {
			return name, n, nil
		}
		if err != nil {
			return "", 0, err
		}

		// If successful, append label to the name.
		name += label + "."
		if len(name) > MaxNameLen {
			return "", 0, ErrNameTooLong
		}
	}
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
	labelLen := endOffset - offset
	if labelLen > MaxLabelLen {
		return "", 0, ErrLabelTooLong
	}
	if !checkBounds(b, offset, endOffset) {
		return "", 0, io.ErrShortBuffer
	}

	if !isPointer {
		n += labelLen
	}
	return string(b[offset:endOffset]), n, nil
}

// Check if begin and end are within bounds of a byte slice.
func checkBounds(b []byte, begin, end int) bool {
	return len(b) >= begin && len(b[begin:]) >= end-begin
}
