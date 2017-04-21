package dns

import (
	"io"
	"regexp"
)

const (
	isPtrOffset = 6

	ptrAddrMask   = 64
	ptrAddrOffset = 8
)

var (
	labelRe = regexp.MustCompile(`^[[:alnum:]][[:alnum:]\-]{0,61}[[:alnum:]]|[[:alpha:]]$^.*[[:^digit:]].*$`)
)

func unpackLabelPointer(b []byte, offset int) (label string, n int, err error) {
	pointerByte := b[offset]
	offset++
	if !checkBounds(b, offset) {
		return "", 0, io.ErrShortBuffer
	}
	pointerOffset := int((pointerByte&ptrAddrMask)<<ptrAddrOffset + b[offset])
	if pointerOffset >= offset-1 {
		return "", 0, ErrLabelPointerIllegal
	}
	label, _, err = unpackLabel(b, pointerOffset)
	return label, 2, err
}

func unpackLabel(b []byte, offset int) (label string, n int, err error) {
	labelLen := int(b[offset])
	if labelLen <= 0 {
		return "", 0, ErrLabelEmpty
	} else if labelLen > MaxLabelLen {
		return "", 0, ErrLabelTooLong
	}

	offset++
	endOffset := offset + labelLen
	if !checkBounds(b, endOffset-1) {
		return "", 0, io.ErrShortBuffer
	}

	label = string(b[offset:endOffset])
	if !isValidLabel(label) {
		return "", 0, ErrLabelInvalid
	}
	return label, labelLen + 1, nil
}

func isPointer(b byte) bool {
	return b>>isPtrOffset == 3
}

func isSafePointer(b []byte, offset int, pointerTable map[int]bool) bool {
	pointerOffset := int((b[offset]&ptrAddrMask)<<ptrAddrMask + b[offset+1])
	_, ok := pointerTable[pointerOffset]
	return ok
}

func isValidLabel(label string) bool {
	return labelRe.MatchString(label)
}
