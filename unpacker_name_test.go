package dns

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	emptyBuffer = bytes.Buffer{}
	eofBuffer   = bytes.Buffer{}
	_           = initBuffers()
)

func initBuffers() bool {
	eofBuffer.WriteByte(0)
	return true
}

// unpackName tests

func TestUnpackNameEmpty(t *testing.T) {
	label, n, err := unpackName(emptyBuffer.Bytes(), 0)
	assert.Error(t, err)
	assert.Equal(t, io.ErrShortBuffer, err)
	assert.Equal(t, "", label)
	assert.Equal(t, 0, n)
}

func TestUnpackNameEOF(t *testing.T) {
	label, n, err := unpackName(eofBuffer.Bytes(), 0)
	assert.NoError(t, err)
	assert.Equal(t, "", label)
	assert.Equal(t, 1, n)
}

func TestUnpackName(t *testing.T) {
	b := bytes.Buffer{}
	b.WriteByte(6)
	b.WriteString("domain")
	b.WriteByte(0)

	name, n, err := unpackName(b.Bytes(), 0)
	assert.NoError(t, err)
	assert.Equal(t, "domain.", name)
	assert.Equal(t, 8, n)
}

// unpackLabel tests

func TestUnpackLabelEmpty(t *testing.T) {
	label, n, err := unpackLabel(emptyBuffer.Bytes(), 0)
	assert.Error(t, err)
	assert.Equal(t, io.ErrShortBuffer, err)
	assert.Equal(t, "", label)
	assert.Equal(t, 0, n)
}

func TestUnpackLabelEOF(t *testing.T) {
	label, n, err := unpackLabel(eofBuffer.Bytes(), 0)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, "", label)
	assert.Equal(t, 1, n)
}

func TestUnpackLabel(t *testing.T) {
	b := bytes.Buffer{}
	b.WriteByte(6)
	b.WriteString("domain")
	b.WriteByte(0)

	label, n, err := unpackLabel(b.Bytes(), 0)
	if assert.NoError(t, err) {
		assert.Equal(t, "domain", label)
		assert.Equal(t, 7, n)
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

	label, n, err := unpackLabel(b.Bytes(), 7)
	if assert.NoError(t, err) {
		assert.Equal(t, "domain", label)
		assert.Equal(t, 1, n)
	}
}
