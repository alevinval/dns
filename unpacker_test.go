package dns

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackMsg(t *testing.T) {
	cases := []struct {
		Input string
		Err   error
	}{
		{Input: "", Err: io.ErrShortBuffer},
		{Input: "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00", Err: io.ErrShortBuffer},

		// Test short buffer by QDCount and ANCount
		{Input: "\x00\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00", Err: io.ErrShortBuffer},
		{Input: "\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x00", Err: io.ErrShortBuffer},

		// NSCount and ARCount still not implemented
		// {Input: "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00", Err: io.ErrShortBuffer},
		// {Input: "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01", Err: io.ErrShortBuffer},

		{Input: "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"},
	}
	for _, c := range cases {
		_, n, err := UnpackMsg([]byte(c.Input), 0)
		assert.Equal(t, c.Err, err)
		if err == nil {
			assert.Equal(t, len(c.Input), n)
		}
	}
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
	b := &bytes.Buffer{}
	packHeader(b, inputHeader)

	unpackedHeader, n, err := unpackHeader(b.Bytes(), 0)
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
