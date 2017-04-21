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
	ErrLabelEmpty          = errors.New("label cannot be empty")
	ErrLabelTooLong        = errors.New("label must be 63 octets or less")
	ErrLabelInvalid        = errors.New("label format is invalid")
	ErrLabelPointerIllegal = errors.New("label pointer is illegal")
	ErrNameTooLong         = errors.New("name must be 255 octets or less")
)

func UnpackMsg(b []byte, offset int) (msg *Msg, n int, err error) {
	initialOffset := offset

	h, n, err := unpackHeader(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	msg = &Msg{Header: *h}

	pointerTable := map[int]bool{}
	msg.Queries = make([]Query, msg.Header.QDCount)
	for i := range msg.Queries {
		q, n, err := unpackQuery(b, offset, pointerTable)
		if err != nil {
			return nil, 0, err
		}
		offset += n
		msg.Queries[i] = *q
	}

	pointerTable = map[int]bool{}
	msg.Responses = make([]RR, msg.Header.ANCount)
	for i := range msg.Responses {
		rr, n, err := unpackRR(b, offset, pointerTable)
		if err != nil {
			return nil, 0, err
		}
		offset += n
		msg.Responses[i] = *rr
	}
	return msg, offset - initialOffset, nil
}

func unpackHeader(b []byte, offset int) (h *Header, n int, err error) {
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

func unpackQuery(b []byte, offset int, pointerTable map[int]bool) (q *Query, n int, err error) {
	initialOffset := offset

	qName, n, err := unpackName(b, offset, pointerTable)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	qtype, n, err := unpackUint16(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	qclass, n, err := unpackUint16(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	q = &Query{QName: qName, QType: Type(qtype), QClass: Class(qclass)}
	return q, offset - initialOffset, nil

}

func unpackRR(b []byte, offset int, pointerTable map[int]bool) (r *RR, n int, err error) {
	initialOffset := offset

	name, n, err := unpackName(b, offset, pointerTable)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	rrtype, n, err := unpackUint16(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	class, n, err := unpackUint16(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	ttl, n, err := unpackUint32(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	rdlength, n, err := unpackUint16(b, offset)
	if err != nil {
		return nil, 0, err
	}
	offset += n

	if len(b[offset:]) < int(rdlength) {
		return nil, 0, io.ErrShortBuffer
	}
	rdata := b[offset : offset+int(rdlength)]
	offset += int(rdlength)

	return &RR{
		Name:     name,
		Class:    Class(class),
		Type:     Type(rrtype),
		TTL:      ttl,
		RDLength: rdlength,
		RData:    rdata,
	}, offset - initialOffset, nil
}

func unpackName(b []byte, offset int, pointerTable map[int]bool) (name string, n int, err error) {
	var ln int
	var label string
	for {
		if !checkBounds(b, offset) {
			return "", 0, io.ErrShortBuffer
		}
		currentByte := b[offset]
		if currentByte == 0 {
			return name, n + 1, nil
		}

		if isPointer(currentByte) {
			if !checkBounds(b, offset+2) {
				return "", 0, io.ErrShortBuffer
			}
			if !isSafePointer(b, offset, pointerTable) {
				return "", 0, ErrLabelPointerIllegal
			}
			label, ln, err = unpackLabelPointer(b, offset)
		} else {
			label, ln, err = unpackLabel(b, offset)
			if err == nil {
				pointerTable[offset] = true
			}
		}
		if err != nil {
			return "", 0, err
		}

		name += label + "."
		offset += ln
		n += ln
		if len(name) > MaxNameLen {
			return "", 0, ErrNameTooLong
		}
	}
}

func unpackUint16(b []byte, offset int) (r uint16, n int, err error) {
	end := offset + 1
	if !checkBounds(b, end) {
		return 0, 0, io.ErrShortBuffer
	}
	return uint16(b[offset])<<8 | uint16(b[end]), 2, nil
}

func unpackUint32(b []byte, offset int) (r uint32, n int, err error) {
	end := offset + 3
	if !checkBounds(b, end) {
		return 0, 0, io.ErrShortBuffer
	}
	return uint32(b[offset])<<24 | uint32(b[offset+1])<<16 | uint32(b[offset+2])<<8 | uint32(b[end]), 4, nil
}

// Check if begin and end are within bounds of a byte slice.
func checkBounds(b []byte, end int) bool {
	return end >= 0 && end < len(b)
}
