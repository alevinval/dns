package dns

import (
	"errors"
	"io"
	"regexp"
)

const (
	MaxLabelLen = 63
)

var (
	ErrLabelEmpty          = errors.New("label cannot be empty")
	ErrLabelTooLong        = errors.New("label must be 63 octets or less")
	ErrLabelInvalid        = errors.New("label format is invalid")
	ErrLabelPointerIllegal = errors.New("label pointer is illegal")

	labelRe   = regexp.MustCompile(`(^[[:alnum:]]$)|(^[[:alnum:]][[:alnum:]\-]{0,61}[[:alnum:]]$)`)
	numericRe = regexp.MustCompile(`^[[:digit:]]*$`)
)

// This function assumes that the pointer is valid. Hence the label is directly
// returned without properly unpacking it.
func unpackLabelPointer(b []byte, offset int) (label string, n int, err error) {
	pointerOffset := getPointerOffset(b, offset)
	if pointerOffset >= offset-1 {
		return "", 0, ErrLabelPointerIllegal
	}
	l := int(b[pointerOffset])
	pointerOffset++
	return string(b[pointerOffset : pointerOffset+l]), 2, err
}

func unpackLabel(b []byte, offset int) (label string, n int, err error) {
	labelLen := int(b[offset])
	if labelLen == 0 {
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
	return b>>0x06 == 0x03
}

func isSafePointer(b []byte, offset int, pointerTable map[int]struct{}) bool {
	pointerOffset := getPointerOffset(b, offset)
	_, ok := pointerTable[pointerOffset]
	return ok
}

func getPointerOffset(b []byte, offset int) int {
	return int(b[offset]&0x3F)<<0x08 + int(b[offset+1])
}

func isValidLabel(label string) bool {
	return labelRe.MatchString(label) && !numericRe.MatchString(label)
}
