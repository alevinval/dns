package dns

import "fmt"

type (
	QClass uint16
	Class  uint16
)

const (
	QClassANY QClass = 255
)

const (
	ClassIN Class = 1 + iota
	ClassCS
	ClassCH
	ClassHS
)

var (
	qClassToString = map[QClass]string{
		QClassANY: "ANY",
	}
	classToString = map[Class]string{
		ClassIN: "IN",
		ClassCS: "CS",
		ClassCH: "CH",
		ClassHS: "HS",
	}
)

func (qc QClass) String() string {
	s, ok := qClassToString[qc]
	if ok {
		return s
	}
	return Class(qc).String()
}

func (c Class) String() string {
	s, ok := classToString[c]
	if ok {
		return s
	}
	return fmt.Sprintf("invalid(%d)", c)

}

func (qc QClass) MarshalText() ([]byte, error) {
	return []byte(qc.String()), nil
}

func (c Class) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}
