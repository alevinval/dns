package dns

import "fmt"

type (
	QType int16
	Type  int16
)

const (
	QTypeAXFR QType = 252 + iota
	QTypeMAILB
	QTypeMAILA
	QTypeALL
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
	qTypeToString = map[QType]string{
		QTypeAXFR:  "AXFR",
		QTypeMAILB: "MAILB",
		QTypeMAILA: "MAILA",
		QTypeALL:   "ALL",
	}

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
	}
)

func (qt QType) String() string {
	s, ok := qTypeToString[qt]
	if ok {
		return s
	}
	return Type(qt).String()
}

func (t Type) String() string {
	s, ok := typeToString[t]
	if ok {
		return s
	}
	return fmt.Sprintf("invalid(%d)", t)
}

func (qt QType) MarshalText() ([]byte, error) {
	return []byte(qt.String()), nil
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
