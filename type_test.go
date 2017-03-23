package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var types = []Type{TypeAXFR, TypeMAILA, TypeMAILB, TypeALL,
	TypeA, TypeNS, TypeMD, TypeMF, TypeCNAME, TypeSOA,
	TypeMB, TypeMG, TypeMR, TypeNULL, TypeWKS, TypePTR,
	TypeHINFO, TypeMINFO, TypeMX, TypeTXT}

func TestTypeToString(t *testing.T) {
	for _, x := range types {
		assert.NotEmpty(t, x.String())
		assert.NotContains(t, x.String(), "invalid")
		marshal, _ := x.MarshalText()
		assert.Equal(t, x.String(), string(marshal))
	}
}
