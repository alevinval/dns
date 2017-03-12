package dns

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	headerBuffer = bytes.Buffer{}
	_            = initialise()
)

func initialise() bool {
	headerBuffer.Write([]byte{
		// ID
		1, 255,
		// TODO: implement tests for the missing flags.
		// FLAGS
		1<<7 + byte(OpCodeSTATUS)<<3, 255,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	)
	return true
}

func TestUnpackHeaderID(t *testing.T) {
	expectedID := uint16(255 + 256)
	h, n, err := unpackHeader(headerBuffer.Bytes(), 0)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, lenHeader, n)
	assert.Equal(t, expectedID, h.ID)
	assert.True(t, h.QR)
	assert.Equal(t, OpCodeSTATUS, h.OpCode)
}
