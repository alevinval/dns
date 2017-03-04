package dns

import (
	"io"
	"strings"
)

type Writer struct {
	msg *Msg
	b   []byte
	i   int
}

func NewWriter(msg *Msg) *Writer {
	return &Writer{msg: msg}
}

func (w *Writer) Bytes() []byte {
	return w.b[:w.i]
}

func (w *Writer) Write(b []byte) (int, error) {
	if w.i == 0 {
		w.writeMsgToBuffer()
	}
	if len(b) < w.i {
		return w.i, io.ErrShortBuffer
	}
	copy(b, w.Bytes())
	return w.i, nil
}

func (w *Writer) WriteTo(dst io.Writer) (int, error) {
	if w.i == 0 {
		w.writeMsgToBuffer()
	}
	return dst.Write(w.Bytes())
}

func (w *Writer) writeMsgToBuffer() {
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
		w.b = append(w.b, byte(labelLen))
		w.b = append(w.b, label...)
		w.i += 1 + labelLen
	}
}

func (w *Writer) writeOctetPair(v uint16) {
	w.b = append(w.b, byte(v>>8))
	w.b = append(w.b, byte(v))
	w.i += 2
}
