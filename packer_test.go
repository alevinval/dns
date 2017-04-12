package dns

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackerUsesPointers(t *testing.T) {
	expectedMsg := &Msg{Header: Header{QDCount: 2}, Queries: []Query{
		{QName: "www.ns1.domain.com.", QType: TypeALL, QClass: ClassANY},
		{QName: "www.ns2.domain.com.", QType: TypeALL, QClass: ClassANY},
	}, Responses: []RR{}}
	b := PackMsg(expectedMsg)
	actualMsg, _, err := UnpackMsg(b, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedMsg, actualMsg)

	// Assert we saved 9 bytes thanks to pointers.
	assert.Equal(t, headerLen+24+24-9, len(b))
}

func TestPackRR(t *testing.T) {
	expectedMsg := &Msg{Header: Header{QDCount: 1, ANCount: 1}, Queries: []Query{
		{QName: "www.ns1.domain.com.", QType: TypeALL, QClass: ClassANY},
	}, Responses: []RR{
		{Name: "www.ns1.domain.com.",
			Type:     TypeALL,
			Class:    ClassANY,
			TTL:      3213,
			RDLength: 4,
			RData:    net.IPv4(1, 1, 1, 1).To4()},
	}}

	b := PackMsg(expectedMsg)

	actualMsg, _, err := UnpackMsg(b, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedMsg, actualMsg)

}
