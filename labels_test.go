package dns

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackLabel(t *testing.T) {
	cases := []struct {
		Input     string
		Expected  string
		Offset    int
		IsPointer bool
		Err       error
	}{
		{Input: "\x00", Err: ErrLabelEmpty},
		{Input: "\x40", Err: ErrLabelTooLong},
		{Input: "\x01", Err: io.ErrShortBuffer},
		{Input: "\x3f", Err: io.ErrShortBuffer},

		{Input: "\x01.", Err: ErrLabelInvalid},
		{Input: "\x02..", Err: ErrLabelInvalid},
		{Input: "\x07.00000A", Err: ErrLabelInvalid},

		{Input: "\x02ok", Expected: "ok"},
		{Input: "\x03123", Expected: "123"},
		{Input: "\x0asome.email", Expected: "some.email"},

		// Pointer tests
		{Input: "\xc0", Err: io.ErrShortBuffer},
		{Input: "\xc0\x02", Err: ErrLabelPointerIllegal},
		{Input: "\x01\xc0\x00", Offset: 1, Err: ErrLabelInvalid},
		{Input: "\x00\xc0\x00", Offset: 1, Err: ErrLabelEmpty},
		{Input: "\x06domain\xc0\x00\x00", Expected: "domain", Offset: 7, IsPointer: true},
	}

	for _, c := range cases {
		t.Logf("Running case with input %q\n", c.Input)
		label, n, err := unpackLabel([]byte(c.Input), c.Offset)
		assert.Equal(t, c.Err, err)
		if err == nil {
			if c.IsPointer {
				assert.Equal(t, 2, n)
			} else {
				assert.Equal(t, len(c.Input), n)
			}
		}
		assert.Equal(t, c.Expected, label)
	}
}
