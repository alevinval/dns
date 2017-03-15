package dns

import "bytes"

func packHeader(h *Header) []byte {
	var flags uint16
	if h.QR {
		flags += 1 << maskQROffset
	}
	flags += uint16(h.OpCode) << maskOpCodeOffset
	if h.AA {
		flags += 1 << maskAAOffset
	}
	if h.TC {
		flags += 1 << maskTCOffset
	}
	if h.RD {
		flags += 1 << maskRDOffset
	}
	if h.RA {
		flags += 1 << maskRAOffset
	}
	flags += uint16(h.Z) << maskZOffset
	flags += uint16(h.RCode)

	header := bytes.NewBuffer(make([]byte, 0, lenHeader))
	header.WriteByte(byte(h.ID >> 8))
	header.WriteByte(byte(h.ID))
	header.WriteByte(byte(flags >> 8))
	header.WriteByte(byte(flags))
	header.WriteByte(byte(h.QDCount >> 8))
	header.WriteByte(byte(h.QDCount))
	header.WriteByte(byte(h.ANCount >> 8))
	header.WriteByte(byte(h.ANCount))
	header.WriteByte(byte(h.NSCount >> 8))
	header.WriteByte(byte(h.NSCount))
	header.WriteByte(byte(h.ARCount >> 8))
	header.WriteByte(byte(h.ARCount))
	return header.Bytes()
}
