package dns

import (
	"io"
)

func unpackName(b []byte, offset int, pointerTable map[int]bool) (name string, n int, err error) {
	var label string
	var ln int
	var initialOffset = offset
	for offset < len(b) {
		switch {
		case b[offset] == 0:
			if len(name) == 0 {
				return "", 0, ErrNameEmpty
			}
			return name, offset - initialOffset + 1, nil
		case isPointer(b[offset]):
			if !checkBounds(b, offset+2) {
				return "", 0, io.ErrShortBuffer
			} else if !isSafePointer(b, offset, pointerTable) {
				return "", 0, ErrLabelPointerIllegal
			}
			label, ln, err = unpackLabelPointer(b, offset)
		default:
			label, ln, err = unpackLabel(b, offset)
			if err == nil {
				pointerTable[offset] = true
			}
		}

		if err != nil {
			return "", 0, err
		} else if len(name)+ln+1 > MaxNameLen {
			return "", 0, ErrNameTooLong
		}

		name += label + "."
		offset += ln
	}
	return "", 0, io.ErrShortBuffer
}
