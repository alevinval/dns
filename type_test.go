package dns

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypePackAndUnpack(t *testing.T) {
	cases := []struct {
		Input    string
		Expected Type
		Err      error
	}{
		{Input: "\x00\x01", Expected: TypeA},
		{Input: "\x00\x02", Expected: TypeNS},
		{Input: "\x00\x03", Expected: TypeMD},
		{Input: "\x00\x04", Expected: TypeMF},
		{Input: "\x00\x05", Expected: TypeCNAME},
		{Input: "\x00\x06", Expected: TypeSOA},
		{Input: "\x00\x07", Expected: TypeMB},
		{Input: "\x00\x08", Expected: TypeMG},
		{Input: "\x00\x09", Expected: TypeMR},
		{Input: "\x00\x0a", Expected: TypeNULL},
		{Input: "\x00\x0b", Expected: TypeWKS},
		{Input: "\x00\x0c", Expected: TypePTR},
		{Input: "\x00\x0d", Expected: TypeHINFO},
		{Input: "\x00\x0e", Expected: TypeMINFO},
		{Input: "\x00\x0f", Expected: TypeMX},
		{Input: "\x00\x10", Expected: TypeTXT},

		{Input: "\x00\xfc", Expected: TypeAXFR},
		{Input: "\x00\xfd", Expected: TypeMAILB},
		{Input: "\x00\xfe", Expected: TypeMAILA},
		{Input: "\x00\xff", Expected: TypeALL},

		{Input: "\x01\x01", Err: ErrTypeInvalid},
	}

	for _, c := range cases {
		t.Logf("Type pack/unpack input: %q\n", c.Input)

		b := &bytes.Buffer{}
		nsType, n, err := unpackType([]byte(c.Input), 0)
		packErr := packType(b, nsType)

		assert.Equal(t, c.Err, err)
		assert.Equal(t, c.Err, packErr)
		if err == nil {
			assert.Equal(t, 2, n)
			assert.Equal(t, c.Expected, nsType)
			assert.NotContains(t, nsType.String(), "invalid")
			assert.Equal(t, c.Input, b.String())
		} else {
			assert.Contains(t, nsType.String(), "invalid")
		}

		marshaledText, _ := nsType.MarshalText()
		assert.Equal(t, nsType.String(), string(marshaledText))
	}
}
