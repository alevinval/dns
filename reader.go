package dns

import (
	"io"
)

type Reader struct {
	i, n               int
	buffer, bufferSwap []byte

	src      io.Reader
	unpacker *Unpacker
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		buffer:     make([]byte, 20),
		bufferSwap: make([]byte, 20),
		unpacker:   NewUnpacker(nil),
		src:        r,
	}
}
func (r *Reader) Read() (msg *Msg, err error) {
	var n int
	for {
		// Read buffer to parse.
		n, err = r.src.Read(r.buffer[r.n:])
		r.n += n
		if err != nil {
			return
		}

		// Unpack as many messages as possible from the available buffer.
		r.unpacker.Reset(r.buffer[r.i:r.n])
		for err != io.ErrShortBuffer {
			msg, n, err = r.unpacker.Unpack()
			if err == io.EOF {
				r.i += n
				return msg, nil
			}
		}
		if r.i == 0 {
			// Short buffer without advancing, buffers are too small.
			r.grow()
		} else {
			// All messages have been serialised, and one is half-way there.
			// Pack buffers to make room for the pending buffer.
			r.pack()
		}
	}
}

func (r *Reader) grow() {
	r.bufferSwap = make([]byte, len(r.buffer)*2)
	data2 := make([]byte, len(r.buffer)*2)
	copy(data2, r.buffer)
	r.buffer = data2
}

func (r *Reader) pack() {
	copy(r.bufferSwap, r.buffer[r.i:r.n])
	copy(r.buffer, r.bufferSwap)
	r.n -= r.i
	r.i = 0
}
