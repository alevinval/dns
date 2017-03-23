package dns

import (
	"errors"
	"io"
)

const (
	MaxNameLen  = 255
	MaxLabelLen = 63

	headerLen = 12
)

var (
	ErrLabelTooLong = errors.New("a label must be 63 octets or less")
	ErrNameTooLong  = errors.New("a name must be 255 octets or less")
)

func UnpackMsg(b []byte, offset int) (msg *Msg, n int, err error) {
	initialOffset := offset

	h, n, err := unpackHeader(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	msg = &Msg{Header: *h}
	msg.Queries = make([]Query, msg.Header.QDCount)

	for i := range msg.Queries {
		q, n, err := unpackQuery(b, offset)
		if err != nil {
			return nil, 0, err
		}
		offset += n
		msg.Queries[i] = *q
	}
	return msg, offset - initialOffset, nil
}

func unpackHeader(b []byte, offset int) (h *Header, n int, err error) {
	if !checkBounds(b, offset+headerLen-1) {
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
	if err != nil {
		return nil, 0, err
	}

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

func unpackQuery(b []byte, offset int) (q *Query, n int, err error) {
	initialOffset := offset

	qName, n, err := unpackName(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	qType, n, err := unpackUint16(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	qClass, n, err := unpackUint16(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	q = &Query{QName: qName, QType: Type(qType), QClass: Class(qClass)}
	return q, offset - initialOffset, nil

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
	if !checkBounds(b, offset) {
		return "", 0, io.ErrShortBuffer
	}

	// Current byte indicates the length of the label.
	// If its a null byte, label is over.
	currentByte := b[offset]
	if currentByte == 0 {
		return "", 1, io.EOF
	}
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
	if !checkBounds(b, endOffset) {
		return "", 0, io.ErrShortBuffer
	}

	if !isPointer {
		n = labelLen + 1
	} else {
		n = 2
	}
	return string(b[offset:endOffset]), n, nil
}

func unpackUint16(b []byte, offset int) (r uint16, n int, err error) {
	end := offset + 1
	if !checkBounds(b, end) {
		return 0, 0, io.ErrShortBuffer
	}
	return uint16(b[end]) | uint16(b[offset])<<8, 2, nil
}

// Check if begin and end are within bounds of a byte slice.
func checkBounds(b []byte, end int) bool {
	return end >= 0 && end < len(b)
}
