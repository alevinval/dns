package dns

import "strings"

type Writer struct {
	msg  *Msg
	data []byte
	i    int
}

// TODO: this is so wip.
func NewWriter(msg *Msg) *Writer {
	return &Writer{msg: msg}
}

func (w *Writer) Bytes() []byte {
	return w.data[:w.i]
}

func (w *Writer) Write() {
	w.writeHeader()
	w.writeQueries()
}

func (w *Writer) writeHeader() {
	header := w.msg.Header
	w.writeOctetPair(header.ID)

	var flags uint16 = 0
	if header.QR {
		flags |= 1 << maskQROffset
	}
	flags |= uint16(header.OpCode) << maskOpCodeOffset
	if header.AA {
		flags |= 1 << maskAAOffset
	}
	if header.TC {
		flags |= 1 << maskTCOffset
	}
	if header.RD {
		flags |= 1 << maskRDOffset
	}
	if header.RA {
		flags |= 1 << maskRAOffset
	}
	flags |= uint16(header.Z) << maskZOffset
	flags |= uint16(header.RCode) << maskRCode
	w.writeOctetPair(flags)

	w.writeOctetPair(header.QDCount)
	w.writeOctetPair(header.ANCount)
	w.writeOctetPair(header.NSCount)
	w.writeOctetPair(header.ARCount)
}

func (w *Writer) writeQueries() {
	for _, query := range w.msg.Queries {
		w.writeName(query.QName)
		w.writeOctetPair(uint16(query.QType))
		w.writeOctetPair(uint16(query.QClass))
	}
}

func (w *Writer) writeName(name string) {
	labels := strings.Split(name, ".")
	for _, label := range labels {
		labelLen := len(label)
		w.data = append(w.data, byte(labelLen))
		w.data = append(w.data, label...)
		w.i += 1 + labelLen
	}
}

func (w *Writer) writeOctetPair(v uint16) {
	w.data = append(w.data, byte(v>>8))
	w.data = append(w.data, byte(v))
	w.i += 2
}
