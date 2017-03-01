package dns

import "strconv"

type (
	QClass uint16
	Class  uint16
)

func (qc QClass) String() string {
	switch qc {
	case 255:
		return "*"
	default:
		return Class(qc).String()
	}
}

func (c Class) String() string {
	switch c {
	case 1:
		return "IN"
	case 2:
		return "CS"
	case 3:
		return "CH"
	case 4:
		return "HS"
	default:
		return "UNKNOWN: " + strconv.Itoa(int(c))
	}
}

func (qc QClass) MarshalText() ([]byte, error) {
	return []byte(qc.String()), nil
}

func (c Class) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}
