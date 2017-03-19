package dns

import (
	"bytes"
	"strings"
)

func PackMsg(msg *Msg) []byte {
	b := &bytes.Buffer{}
	PackMsgTo(b, msg)
	return b.Bytes()
}

func PackMsgTo(b *bytes.Buffer, msg *Msg) {
	packHeader(b, &msg.Header)
	for i := range msg.Queries {
		packQuery(b, &msg.Queries[i])
	}
}

func packHeader(b *bytes.Buffer, h *Header) {
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

	writeUint16(b, h.ID)
	writeUint16(b, flags)
	writeUint16(b, h.QDCount)
	writeUint16(b, h.ANCount)
	writeUint16(b, h.NSCount)
	writeUint16(b, h.ARCount)
}

func packQuery(b *bytes.Buffer, q *Query) {
	writeName(b, q.QName)
	writeUint16(b, uint16(q.QType))
	writeUint16(b, uint16(q.QClass))
}

func writeName(b *bytes.Buffer, name string) {
	labels := strings.Split(name, ".")
	for _, label := range labels {
		labelLen := len(label)
		b.WriteByte(byte(labelLen))
		b.WriteString(label)
	}
}

func writeUint16(b *bytes.Buffer, v uint16) {
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v))
}
