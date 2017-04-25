package dns

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	TypeAXFR Type = 252 + iota
	TypeMAILB
	TypeMAILA
	TypeALL
)

const (
	TypeA Type = 1 + iota
	TypeNS
	TypeMD
	TypeMF
	TypeCNAME
	TypeSOA
	TypeMB
	TypeMG
	TypeMR
	TypeNULL
	TypeWKS
	TypePTR
	TypeHINFO
	TypeMINFO
	TypeMX
	TypeTXT
)

var (
	ErrTypeInvalid = errors.New("invalid type value")
)

var (
	typeToString = map[Type]string{
		TypeA:     "A",
		TypeNS:    "NS",
		TypeMD:    "MD",
		TypeMF:    "MF",
		TypeCNAME: "CNAME",
		TypeSOA:   "SOA",
		TypeMB:    "MB",
		TypeMG:    "MG",
		TypeMR:    "MR",
		TypeNULL:  "NULL",
		TypeWKS:   "WKS",
		TypePTR:   "PTR",
		TypeHINFO: "HINFO",
		TypeMINFO: "MINFO",
		TypeMX:    "MX",
		TypeTXT:   "TXT",

		TypeAXFR:  "AXFR",
		TypeMAILB: "MAILB",
		TypeMAILA: "MAILA",
		TypeALL:   "ALL",
	}
)

type Type int16

func packType(b *bytes.Buffer, t Type) (err error) {
	_, ok := typeToString[t]
	if !ok {
		return ErrTypeInvalid
	}
	writeUint16(b, uint16(t))
	return
}

func unpackType(b []byte, offset int) (t Type, n int, err error) {
	u, n, err := unpackUint16(b, offset)
	t = Type(u)
	_, ok := typeToString[t]
	if !ok {
		return 0, 0, ErrTypeInvalid
	}
	return t, n, err
}

func (t Type) String() string {
	s, ok := typeToString[t]
	if ok {
		return s
	}
	return fmt.Sprintf("invalid(%d)", t)
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
