package dns

import (
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

		{Input: "\x01\x99", Err: ErrLabelInvalid},
		{Input: "\x3fabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789x\xc0\x00\xc0\x00\xc0\x00\xc0\x00\x00", Err: ErrNameTooLong},

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
