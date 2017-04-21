package dns

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	emptyPayload = []byte{}
	eofBuffer    = []byte{0}
)

func TestUnpackMsgShortBuffer(t *testing.T) {
	_, n, err := UnpackMsg(emptyPayload, 0)
	assert.Error(t, err)
	assert.Equal(t, io.ErrShortBuffer, err)
	assert.Equal(t, 0, n)
}

func TestUnpackMsgEOF(t *testing.T) {
	emptyMessage := make([]byte, 32)
	_, n, err := UnpackMsg(emptyMessage, 0)
	assert.NoError(t, err)
	assert.Equal(t, headerLen, n)
}

func TestUnpackMsg(t *testing.T) {
	expected := &Msg{Header: Header{QDCount: 1}, Queries: []Query{
		{QName: "www.test.com.", QType: TypeALL, QClass: Class(ClassIN)},
	}, Responses: []RR{}}
	b := PackMsg(expected)

	actual, n, err := UnpackMsg(b, 0)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, len(b), n) {
		return
	}
	assert.Equal(t, expected, actual)
}

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

func TestUnpackNameEmpty(t *testing.T) {
	pointerTable := map[int]bool{}
	label, n, err := unpackName(emptyPayload, 0, pointerTable)
	assert.Error(t, err)
	assert.Equal(t, io.ErrShortBuffer, err)
	assert.Equal(t, "", label)
	assert.Equal(t, 0, n)
}

func TestUnpackNameEOF(t *testing.T) {
	pointerTable := map[int]bool{}
	label, n, err := unpackName(eofBuffer, 0, pointerTable)
	assert.NoError(t, err)
	assert.Equal(t, "", label)
	assert.Equal(t, 1, n)
}

func TestUnpackNameTooLong(t *testing.T) {
	maxLabel := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789x"

	b := bytes.Buffer{}
	for i := 0; i < 4; i++ {
		b.WriteByte(63)
		b.WriteString(maxLabel)
	}
	b.WriteByte(0)

	pointerTable := map[int]bool{}
	_, _, err := unpackName(b.Bytes(), 0, pointerTable)
	assert.Equal(t, ErrNameTooLong, err)
}

func TestUnpackName(t *testing.T) {
	b := bytes.Buffer{}
	b.WriteByte(6)
	b.WriteString("domain")
	b.WriteByte(0)

	pointerTable := map[int]bool{}
	name, n, err := unpackName(b.Bytes(), 0, pointerTable)
	assert.NoError(t, err)
	assert.Equal(t, "domain.", name)
	assert.Equal(t, 8, n)
}
