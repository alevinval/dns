package dns

import (
	"bytes"
	"io"
	"strings"
)

const MaxNameLen = 255

func packName(b *bytes.Buffer, labelTable map[string]int, name string) error {
	name = strings.TrimSuffix(name, ".")
	labels := strings.Split(name, ".")
	for _, label := range labels {
		if !isValidLabel(label) {
			return ErrLabelInvalid
		}
		position, seen := labelTable[label]
		if seen {
			packPointerTo(b, position)
		} else {
			l := len(label)
			// No point in using pointers for labels of length 1.
			if l > 1 {
				labelTable[label] = b.Len()
			}
			b.WriteByte(byte(l))
			b.WriteString(label)
		}
	}
	b.WriteByte(0)
	return nil
}

func unpackName(b []byte, offset int, pointerTable map[int]struct{}) (name string, n int, err error) {
	var label string
	var ln int
	initialOffset := offset
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
				pointerTable[offset] = struct{}{}
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

func packPointerTo(b *bytes.Buffer, offset int) {
	b.WriteByte(byte(offset>>8 | 3<<6))
	b.WriteByte(byte(offset))
}
