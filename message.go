package dns

type (
	Header struct {
		ID      uint16
		Flags   uint16
		QDCount uint16
		ANCount uint16
		NSCount uint16
		ARCount uint16
	}
	Query struct {
		QName  string
		QType  QType
		QClass QClass
	}
	Msg struct {
		Header  Header
		Queries []Query
	}
)
