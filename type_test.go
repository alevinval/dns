package dns

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var qTypes = []QType{QTypeAXFR, QTypeMAILA, QTypeMAILB, QTypeALL}
var types = []Type{TypeA, TypeNS, TypeMD, TypeMF, TypeCNAME, TypeSOA,
	TypeMB, TypeMG, TypeMR, TypeNULL, TypeWKS, TypePTR,
	TypeHINFO, TypeMINFO, TypeMX, TypeTXT}

func TestQTypeToString(t *testing.T) {
	for _, qType := range qTypes {
		assert.NotEmpty(t, qType.String())
		assert.NotContains(t, qType.String(), "invalid")
		marshal, _ := qType.MarshalText()
		assert.Equal(t, qType.String(), string(marshal))
	}
}

func TestTypeToString(t *testing.T) {
	for _, x := range types {
		assert.NotEmpty(t, x.String())
		assert.NotContains(t, x.String(), "invalid")
		marshal, _ := x.MarshalText()
		assert.Equal(t, x.String(), string(marshal))
	}
}

func TestQTypeIsSuperset(t *testing.T) {
	for _, tType := range types {
		qType := QType(tType)
		assert.Equal(t, tType.String(), qType.String())
	}
}

func TestTypeIsNotSuperset(t *testing.T) {
	for _, qType := range qTypes {
		tType := Type(qType)
		assert.NotEqual(t, qType.String(), tType.String())
		assert.Equal(t, tType.String(), fmt.Sprintf("invalid(%d)", tType))
	}
}
