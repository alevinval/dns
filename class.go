package dns

import "fmt"

type Class uint16

const (
	ClassANY Class = 255
)

const (
	ClassIN Class = 1 + iota
	ClassCS
	ClassCH
	ClassHS
)

var (
	classToString = map[Class]string{
		// QClass
		ClassANY: "ANY",

		// Class
		ClassIN: "IN",
		ClassCS: "CS",
		ClassCH: "CH",
		ClassHS: "HS",
	}
)

func (c Class) String() string {
	s, ok := classToString[c]
	if ok {
		return s
	}
	return fmt.Sprintf("invalid(%d)", c)

}

func (c Class) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}
