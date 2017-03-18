package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackHeaderShortBuffer(t *testing.T) {
	b := make([]byte, headerLen-1)
	_, n, err := unpackHeader(b, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, n)

	b = make([]byte, 0, 1)
	_, n, err = unpackHeader(b, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, n)
}

func TestUnpackHeader(t *testing.T) {
	inputHeader := &Header{
		ID:      255 + 256,
		QR:      true,
		OpCode:  OpCodeSTATUS,
		AA:      true,
		TC:      true,
		RD:      true,
		RA:      true,
		Z:       6,
		RCode:   5,
		QDCount: 4,
		ANCount: 3,
		NSCount: 2,
		ARCount: 1,
	}
	data := packHeader(inputHeader)

	unpackedHeader, n, err := unpackHeader(data, 0)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, headerLen, n)
	assert.Equal(t, inputHeader.ID, unpackedHeader.ID)
	assert.Equal(t, inputHeader.QR, unpackedHeader.QR)
	assert.Equal(t, inputHeader.OpCode, unpackedHeader.OpCode)
	assert.Equal(t, inputHeader.AA, unpackedHeader.AA)
	assert.Equal(t, inputHeader.TC, unpackedHeader.TC)
	assert.Equal(t, inputHeader.RD, unpackedHeader.RD)
	assert.Equal(t, inputHeader.RA, unpackedHeader.RA)
	assert.Equal(t, inputHeader.Z, unpackedHeader.Z)
	assert.Equal(t, inputHeader.RCode, unpackedHeader.RCode)
	assert.Equal(t, inputHeader.QDCount, unpackedHeader.QDCount)
	assert.Equal(t, inputHeader.ANCount, unpackedHeader.ANCount)
	assert.Equal(t, inputHeader.NSCount, unpackedHeader.NSCount)
	assert.Equal(t, inputHeader.ARCount, unpackedHeader.ARCount)
}
