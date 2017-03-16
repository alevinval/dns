package dns

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var qClasses = []QClass{QClassANY}
var classes = []Class{ClassCH, ClassCS, ClassHS, ClassIN}

func TestClassToString(t *testing.T) {
	for _, class := range classes {
		assert.NotEmpty(t, class.String())
		assert.NotContains(t, class.String(), "invalid")
		marshal, _ := class.MarshalText()
		assert.Equal(t, class.String(), string(marshal))
	}
}

func TestQClassToString(t *testing.T) {
	for _, qClass := range qClasses {
		assert.NotEmpty(t, qClass.String())
		assert.NotContains(t, qClass.String(), "invalid")
		marshal, _ := qClass.MarshalText()
		assert.Equal(t, qClass.String(), string(marshal))
	}
}

func TestQClassIsSuperset(t *testing.T) {
	for _, class := range classes {
		qClass := QClass(class)
		assert.Equal(t, class.String(), qClass.String())
	}
}

func TestClassIsNotSuperset(t *testing.T) {
	for _, qClass := range qClasses {
		class := Class(qClass)
		assert.NotEqual(t, qClass.String(), class.String())
		assert.Equal(t, class.String(), fmt.Sprintf("invalid(%d)", class))
	}
}
