package dns

import "fmt"

type OpCode byte

const (
	OpCodeQUERY = iota
	OpCodeIQUERY
	OpCodeSTATUS
)

var opCodeToString = map[OpCode]string{
	OpCodeQUERY:  "QUERY",
	OpCodeIQUERY: "IQUERY",
	OpCodeSTATUS: "STATUS",
}

func (o OpCode) String() string {
	s, ok := opCodeToString[o]
	if ok {
		return s
	}
	return fmt.Sprintf("invalid(%d)", o)
}

func (o OpCode) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}
