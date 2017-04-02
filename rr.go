package dns

type RR struct {
	Name     string
	Type     Type
	Class    Class
	TTL      uint32
	RDLength uint16
	RData    []byte
}
