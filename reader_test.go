package dns

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderEmpty(t *testing.T) {
	b := make([]byte, 0)
	buff := bytes.NewBuffer(b)
	r := NewReader(buff)

	_, err := r.Read()
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func TestReaderShortBufferNoProgress(t *testing.T) {
	buff := bytes.NewBuffer(make([]byte, 10))
	r := NewReader(buff)

	_, err := r.Read()
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func TestReaderHeaderOnly(t *testing.T) {
	h := &Header{ID: 123}
	b := packHeader(h)
	buff := bytes.NewBuffer(b)
	r := NewReader(buff)

	msg, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, uint16(123), msg.Header.ID)

	_, err = r.Read()
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}
