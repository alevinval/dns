package dns

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackNameSuite(t *testing.T) {
	cases := []struct {
		Input        string
		Expected     string
		PointerTable map[int]bool
		Err          error
	}{
		{Input: "", Err: io.ErrShortBuffer},
		{Input: "\x00", Err: ErrNameEmpty},
		{Input: "\x01a\xc0", Err: io.ErrShortBuffer},
		{Input: "\x01a\xc0\x0f\x00", Err: ErrLabelPointerIllegal},

		{Input: "\x01a\x00", Expected: "a."},
		{Input: "\x01a\xc0\x00\x00", Expected: "a.a."},
	}
	for _, c := range cases {
		t.Logf("Name unpacking input: %q\n", c.Input)
		b := []byte(c.Input)
		if c.PointerTable == nil {
			c.PointerTable = map[int]bool{}
		}
		name, n, err := unpackName(b, 0, c.PointerTable)
		assert.Equal(t, c.Err, err)
		if err == nil {
			assert.Equal(t, c.Expected, name)
			assert.Equal(t, len(c.Input), n)
		} else {
			assert.Equal(t, "", name)
			assert.Equal(t, 0, n)
		}
	}
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
