package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpCodeToString(t *testing.T) {
	opCodes := []OpCode{OpCodeSTATUS, OpCodeQUERY, OpCodeIQUERY}
	for _, opCode := range opCodes {
		assert.NotEmpty(t, opCode.String())
		assert.NotContains(t, opCode.String(), "invalid")
		marshal, _ := opCode.MarshalText()
		assert.Equal(t, opCode.String(), string(marshal))
	}
}

func TestOpCodeInvalid(t *testing.T) {
	o := OpCode(200)
	assert.Equal(t, "invalid(200)", o.String())
}
