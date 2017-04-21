package dns

import (
	"io"
	"regexp"
)

var (
	labelRe = regexp.MustCompile(`^[[:alnum:]][[:alnum:]\-]{0,61}[[:alnum:]]|[[:alpha:]]$^.*[[:^digit:]].*$`)
)

func unpackLabel(b []byte, offset int) (label string, n int, err error) {
	// No need to check bounds because the caller, unpackName, already
	// checked the byte at offset.
	currentByte := b[offset]
	offset++

	isPointer := isPointer(currentByte)
	if isPointer {
		if !checkBounds(b, offset) {
			return "", 0, io.ErrShortBuffer
		}

		// Compute the offset to the pointer.
		pointerOffset := int((currentByte&64)<<8 + b[offset])
		if pointerOffset >= offset {
			return "", 0, ErrLabelPointerIllegal
		}

		currentByte = b[pointerOffset]
		offset = pointerOffset + 1
		n = 2
	}

	endOffset := offset + int(currentByte)

	// Check if the label has valid length.
	labelLen := endOffset - offset
	if labelLen <= 0 {
		return "", 0, ErrLabelEmpty
	} else if labelLen > MaxLabelLen {
		return "", 0, ErrLabelTooLong
	} else if !checkBounds(b, endOffset-1) {
		return "", 0, io.ErrShortBuffer
	}

	label = string(b[offset:endOffset])
	if !isValidLabel(label) {
		return "", 0, ErrLabelInvalid
	}

	if !isPointer {
		n = labelLen + 1
	}
	return label, n, nil
}

func isPointer(b byte) bool {
	return b>>6 == 3
}

func isValidLabel(label string) bool {
	return labelRe.MatchString(label)
}
