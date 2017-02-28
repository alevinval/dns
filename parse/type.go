package parse

import "strconv"

type (
	QType int16
	Type  int16
)

func (qt QType) String() string {
	switch qt {
	case 252:
		return "AXFR"
	case 253:
		return "MAILB"
	case 254:
		return "MAILA"
	case 255:
		return "*"
	default:
		return Type(qt).String()
	}
}

func (t Type) String() string {
	switch t {
	case 1:
		return "A"
	case 2:
		return "NS"
	case 3:
		return "MD"
	case 4:
		return "MF"
	case 5:
		return "CNAME"
	case 6:
		return "SOA"
	case 7:
		return "MB"
	case 8:
		return "MG"
	case 9:
		return "MR"
	case 10:
		return "NULL"
	case 11:
		return "WKS"
	case 12:
		return "PTR"
	case 13:
		return "HINFO"
	case 14:
		return "MINFO"
	case 15:
		return "MX"
	case 16:
		return "TXT"
	default:
		return "UNKNOWN: " + strconv.Itoa(int(t))
	}
}

func (qt QType) MarshalText() ([]byte, error) {
	return []byte(qt.String()), nil
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
