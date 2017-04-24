package dns

import "bytes"

func PackMsg(msg *Msg) []byte {
	b := &bytes.Buffer{}
	PackMsgTo(b, msg)
	return b.Bytes()
}

func PackMsgTo(b *bytes.Buffer, msg *Msg) {
	packHeader(b, &msg.Header)

	labelTable := map[string]int{}
	for i := 0; i < int(msg.Header.QDCount); i++ {
		packQuery(b, labelTable, &msg.Queries[i])
	}

	labelTable = map[string]int{}
	for i := 0; i < int(msg.Header.ANCount); i++ {
		packRR(b, labelTable, &msg.Responses[i])
	}
}

func packRR(b *bytes.Buffer, labelTable map[string]int, rr *RR) {
	packName(b, labelTable, rr.Name)
	writeUint16(b, uint16(rr.Type))
	packClass(b, rr.Class)
	writeUint32(b, rr.TTL)
	writeUint16(b, rr.RDLength)
	b.Write(rr.RData)
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

func packQuery(b *bytes.Buffer, labelTable map[string]int, q *Query) {
	packName(b, labelTable, q.QName)
	writeUint16(b, uint16(q.QType))
	writeUint16(b, uint16(q.QClass))
}

func writeUint16(b *bytes.Buffer, v uint16) {
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v))
}

func writeUint32(b *bytes.Buffer, v uint32) {
	b.WriteByte(byte(v >> 24))
	b.WriteByte(byte(v >> 16))
	b.WriteByte(byte(v >> 8))
	b.WriteByte(byte(v))
}
