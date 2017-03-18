package dns

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackMsgShortBuffer(t *testing.T) {
	_, n, err := UnpackMsg(emptyBuffer.Bytes(), 0)
	assert.Error(t, err)
	assert.Equal(t, io.ErrShortBuffer, err)
	assert.Equal(t, 0, n)
}

func TestUnpackMsgEOF(t *testing.T) {
	emptyMessage := make([]byte, 32)
	_, n, err := UnpackMsg(emptyMessage, 0)
	assert.NoError(t, err)
	assert.Equal(t, headerLen, n)
}
