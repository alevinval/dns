package debug

import (
	"net"

	"github.com/go-rfc/dns"
)

func FakeResponseMsg(m *dns.Msg) *dns.Msg {
	m.Header.QR = true
	m.Header.RD = false

	for i := 0; i < int(m.Header.QDCount); i++ {
		q := m.Queries[i]
		rr := dns.RR{
			Name:     q.QName,
			Class:    q.QClass,
			Type:     q.QType,
			TTL:      3600,
			RDLength: 4,
			RData:    net.IPv4(1, 2, 3, 4).To4()}
		m.Responses = append(m.Responses, rr)
		m.Header.ANCount++
	}

	return m
}
