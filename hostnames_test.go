package dns

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackName(t *testing.T) {
	cases := []struct {
		Input      string
		Expected   string
		LabelTable map[string]int
		Err        error
	}{
		{Input: "", Err: ErrLabelInvalid},
		{Input: ".", Err: ErrLabelInvalid},
		{Input: "-", Err: ErrLabelInvalid},
		{Input: "-ab", Err: ErrLabelInvalid},
		{Input: "ab-", Err: ErrLabelInvalid},
		{Input: "1", Err: ErrLabelInvalid},
		{Input: "123", Err: ErrLabelInvalid},

		{Input: "a", Expected: "\x01a\x00"},
		{Input: "1a", Expected: "\x021a\x00"},
		{Input: "a1", Expected: "\x02a1\x00"},
		{Input: "abc", Expected: "\x03abc\x00"},
		{Input: "a-c", Expected: "\x03a-c\x00"},
		{Input: "a.b.c", Expected: "\x01a\x01b\x01c\x00"},

		// No pointers for 1 octet labels.
		{Input: "a.a.a", Expected: "\x01a\x01a\x01a\x00"},
		// Pointers otherwise.
		{Input: "ab.ab", Expected: "\x02ab\xc0\x00\x00"},
	}

	for _, c := range cases {
		t.Logf("Name packing input: %q\n", c.Input)
		b := &bytes.Buffer{}
		if c.LabelTable == nil {
			c.LabelTable = map[string]int{}
		}
		err := packName(b, c.LabelTable, c.Input)
		assert.Equal(t, c.Err, err)
		if err == nil {
			assert.Equal(t, c.Expected, b.String())
		}
	}
}

func TestUnpackName(t *testing.T) {
	cases := []struct {
		Input        string
		Expected     string
		PointerTable map[int]struct{}
		Err          error
	}{
		{Input: "", Err: io.ErrShortBuffer},
		{Input: "\x00", Err: ErrNameEmpty},
		{Input: "\x01a\xc0", Err: io.ErrShortBuffer},
		{Input: "\x01a\xc0\x0f\x00", Err: ErrLabelPointerIllegal},

		{Input: "\x01\x99", Err: ErrLabelInvalid},
		{Input: "\x3fabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789x\xc0\x00\xc0\x00\xc0\x00\xc0\x00\x00", Err: ErrNameTooLong},

		{Input: "\x01a\x00", Expected: "a."},
		{Input: "\x01a\xc0\x00\x00", Expected: "a.a."},
	}
	for _, c := range cases {
		t.Logf("Name unpacking input: %q\n", c.Input)
		b := []byte(c.Input)
		if c.PointerTable == nil {
			c.PointerTable = map[int]struct{}{}
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
