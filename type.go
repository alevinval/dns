package dns

import "fmt"

type (
	Type int16
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
	typeToString = map[Type]string{
		// QTypes
		TypeAXFR:  "AXFR",
		TypeMAILB: "MAILB",
		TypeMAILA: "MAILA",
		TypeALL:   "ALL",

		// Types
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
	}
)

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
