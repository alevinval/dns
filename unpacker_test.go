package dns

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackLabelEmpty(t *testing.T) {
	b := bytes.Buffer{}

	label, err := unpackLabel(b.Bytes(), 0)
	assert.Error(t, err)
	assert.Equal(t, io.ErrShortBuffer, err)
	assert.Equal(t, "", label)
}

func TestUnpackLabelEOF(t *testing.T) {
	b := bytes.Buffer{}
	b.WriteByte(0)

	label, err := unpackLabel(b.Bytes(), 0)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, "", label)
}

func TestUnpackLabel(t *testing.T) {
	b := bytes.Buffer{}
	b.WriteByte(6)
	b.WriteString("domain")
	b.WriteByte(0)

	label, err := unpackLabel(b.Bytes(), 0)
	if assert.NoError(t, err) {
		assert.Equal(t, "domain", label)
	}
}

func TestUnpackLabelWithPointer(t *testing.T) {
	b := bytes.Buffer{}
	b.WriteByte(6)
	b.WriteString("domain")

	ptr := byte(3 << 6)
	b.WriteByte(ptr)
	b.WriteByte(0)
	b.WriteByte(0)

	label, err := unpackLabel(b.Bytes(), 7)
	if assert.NoError(t, err) {
		assert.Equal(t, "domain", label)
	}
}
