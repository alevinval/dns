package dns

import (
	"bytes"
	"errors"
	"fmt"
)

type Class uint16

const (
	ClassIN Class = 1 + iota
	ClassCS
	ClassCH
	ClassHS

	ClassANY = 255
)

var (
	ErrClassInvalid = errors.New("invalid class value")

	classToString = map[Class]string{
		ClassIN:  "IN",
		ClassCS:  "CS",
		ClassCH:  "CH",
		ClassHS:  "HS",
		ClassANY: "ANY",
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

func packClass(b *bytes.Buffer, c Class) (err error) {
	_, ok := classToString[c]
	if !ok {
		return ErrClassInvalid
	}
	writeUint16(b, uint16(c))
	return
}

func unpackClass(b []byte, offset int) (c Class, n int, err error) {
	u, n, err := unpackUint16(b, offset)
	c = Class(u)
	_, ok := classToString[c]
	if !ok {
		return 0, 0, ErrClassInvalid
	}
	return c, n, err
}
