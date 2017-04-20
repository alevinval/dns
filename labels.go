package dns

import (
	"io"
	"regexp"
)

var (
	labelRe = regexp.MustCompile(`^[[:alnum:]][[:alnum:]\-]{0,61}[[:alnum:]]|[[:alpha:]]$^.*[[:^digit:]].*$`)
)

func unpackLabel(b []byte, offset int) (label string, n int, err error) {
	if !checkBounds(b, offset) {
		return "", 0, io.ErrShortBuffer
	}

	// Current byte indicates the length of the label.
	// If its a null byte, label is over.
	currentByte := b[offset]
	if currentByte == 0 {
		return "", 1, io.EOF
	}
	offset++

	// Check if its a pointer.
	isPointer := currentByte>>6 == 3
	if isPointer {
		if !checkBounds(b, offset) {
			return "", 0, io.ErrShortBuffer
		}

		// Compute the offset to the pointer.
		originalOffset := offset
		offset = int((currentByte&64)<<8 + b[offset])
		if offset >= originalOffset {
			return "", 0, ErrLabelPointerIllegal
		}

		currentByte = b[offset]
		offset++
	}

	endOffset := offset + int(currentByte)

	// Check if the label has valid length.
	labelLen := endOffset - offset
	if labelLen <= 0 {
		return "", 0, ErrLabelEmpty
	} else if labelLen > MaxLabelLen {
		return "", 0, ErrLabelTooLong
	}

	// Check if the label fits in buffer.
	if !checkBounds(b, endOffset-1) {
		return "", 0, io.ErrShortBuffer
	}

	label = string(b[offset:endOffset])
	if !isValidLabel(label) {
		return "", 0, ErrLabelInvalid
	}

	if !isPointer {
		n = labelLen + 1
	} else {
		n = 2
	}
	return label, n, nil
}

func isValidLabel(label string) bool {
	return labelRe.MatchString(label)
}
