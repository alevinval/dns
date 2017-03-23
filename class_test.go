package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var classes = []Class{ClassANY, ClassCH, ClassCS, ClassHS, ClassIN}

func TestClassToString(t *testing.T) {
	for _, class := range classes {
		assert.NotEmpty(t, class.String())
		assert.NotContains(t, class.String(), "invalid")
		marshal, _ := class.MarshalText()
		assert.Equal(t, class.String(), string(marshal))
	}
}
