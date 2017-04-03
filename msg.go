package dns

type (
	Header struct {
		ID      uint16
		QR      bool
		OpCode  OpCode
		AA      bool
		TC      bool
		RD      bool
		RA      bool
		Z       byte
		RCode   byte
		QDCount uint16
		ANCount uint16
		NSCount uint16
		ARCount uint16
	}
	Query struct {
		QName  string
		QType  Type
		QClass Class
	}
	RR struct {
		Name     string
		Type     Type
		Class    Class
		TTL      uint32
		RDLength uint16
		RData    []byte
	}
	Msg struct {
		Header    Header
		Queries   []Query
		Responses []RR
	}
)
