package dns

import "io"

type Reader struct {
	last, i, n int
	buffer     []byte

	src io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		// Buffer should be bigger. Making it tiny while developing
		// the library to help catch logic flaws.
		buffer: make([]byte, 1),
		src:    r,
	}
}
func (r *Reader) Read() (msg *Msg, err error) {
	var n int
	for {
		n, err = r.src.Read(r.buffer[r.n:])
		r.n += n
		if err != nil && (r.i >= r.n || r.last >= r.i) {
			return
		}

		// Unpack message if possible.
		msg, n, err = UnpackMsg(r.buffer[r.i:r.n], 0)
		r.last = r.i
		r.i += n
		if err == io.ErrShortBuffer {
			// Only grow the buffer when really needed.
			if r.i > 0 && r.i > len(r.buffer)/2 {
				r.pack()
			} else {
				r.grow()
			}
			continue
		}
		return
	}
}

func (r *Reader) grow() {
	bigger := make([]byte, len(r.buffer)*2)
	copy(bigger, r.buffer)
	r.buffer = bigger
}

func (r *Reader) pack() {
	copy(r.buffer, r.buffer[r.i:r.n])
	r.n -= r.i
	r.i = 0
}
