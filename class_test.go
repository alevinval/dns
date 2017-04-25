package dns

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassPackAndUnpack(t *testing.T) {
	cases := []struct {
		Input    string
		Expected Class
		Err      error
	}{
		{Input: "\x00\x01", Expected: ClassIN},
		{Input: "\x00\x02", Expected: ClassCS},
		{Input: "\x00\x03", Expected: ClassCH},
		{Input: "\x00\x04", Expected: ClassHS},

		{Input: "\x00\xff", Expected: ClassANY},

		{Input: "\x01\x00", Err: ErrClassInvalid},
	}

	for _, c := range cases {
		t.Logf("Class pack/unpack input: %q\n", c.Input)

		b := &bytes.Buffer{}
		class, n, err := unpackClass([]byte(c.Input), 0)
		packErr := packClass(b, class)

		assert.Equal(t, c.Err, err)
		assert.Equal(t, c.Err, packErr)
		if err == nil {
			assert.Equal(t, 2, n)
			assert.Equal(t, c.Expected, class)
			assert.NotContains(t, class.String(), "invalid")
			assert.Equal(t, c.Input, b.String())
		} else {
			assert.Contains(t, class.String(), "invalid")
		}

		marshaledText, _ := class.MarshalText()
		assert.Equal(t, class.String(), string(marshaledText))
	}
}
