package dns

import "bytes"

func packHeader(h *Header) []byte {
	var flags uint16
	if h.QR {
		flags |= 1 << maskQROffset
	}
	flags |= uint16(h.OpCode) << maskOpCodeOffset
	if h.AA {
		flags |= 1 << maskAAOffset
	}
	if h.TC {
		flags |= 1 << maskTCOffset
	}
	if h.RD {
		flags |= 1 << maskRDOffset
	}
	if h.RA {
		flags |= 1 << maskRAOffset
	}
	flags |= uint16(h.Z) << maskZOffset
	flags |= uint16(h.RCode)

	buffer := bytes.NewBuffer(make([]byte, 0, headerLen))
	writeUint16(buffer, h.ID)
	writeUint16(buffer, flags)
	writeUint16(buffer, h.QDCount)
	writeUint16(buffer, h.ANCount)
	writeUint16(buffer, h.NSCount)
	writeUint16(buffer, h.ARCount)
	return buffer.Bytes()
}

func writeUint16(b *bytes.Buffer, v uint16) {
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v))
}
