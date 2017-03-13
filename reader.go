package dns

import "io"

type Reader struct {
	i, n               int
	buffer, bufferSwap []byte

	src io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		buffer:     make([]byte, 20),
		bufferSwap: make([]byte, 20),
		src:        r,
	}
}
func (r *Reader) Read() (msg *Msg, err error) {
	var n int
	for {
		// Read buffer to parse.
		n, err = r.src.Read(r.buffer[r.n:])
		r.n += n
		if r.i == r.n && err == io.EOF {
			return
		}

		// Unpack message if possible.
		msg, n, err = UnpackMsg(r.buffer[r.i:r.n], 0)
		if err == nil {
			r.i += n
			return msg, nil
		} else if err == io.ErrShortBuffer {
			r.grow()
		} else if r.i > 0 {
			r.pack()
		} else {
			return
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
