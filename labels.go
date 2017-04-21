package dns

import (
	"io"
	"regexp"
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
	pointerOffset := int((pointerByte&64)<<8 + b[offset])
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
	return b>>6 == 3
}

func isValidLabel(label string) bool {
	return labelRe.MatchString(label)
}
