package utils

import (
	"strings"
)

func ParseStringNUL(data []byte) (string, int) {
	var s []byte
	for i, b := range data {
		if b == 0x00 {
			return string(s), i + 1
		} else if IsPrintable(b) {
			s = append(s, b)
		}
	}

	return "", -1
}

func IsPrintable(b byte) bool {
	return 0x20 < b && b < 0x7f
}

func ParseFixUint(data []byte) (num uint64) {
	for i, b := range data {
		num |= uint64(b) << (8 * uint(i))
	}
	return
}

func ParseLenEncUint(data []byte) (num uint64, byte_cnt int) {
	switch data[0] {
	default:
		num = ParseFixUint(data[0:1])
		byte_cnt = 1
		return
	case 0xfb:
		num = 256
		byte_cnt = 1
		return
	case 0xfc:
		byte_cnt = 3
	case 0xfd:
		byte_cnt = 4
	case 0xfe:
		byte_cnt = 9
	}
	num = ParseFixUint(data[1:byte_cnt])
	return

}

func ParseStringNULEOF(data []byte) (string, int) {
	var offset, next int
	var s []string

	total_len := len(data)
	for offset < total_len {
		v, dlen := ParseStringNUL(data[offset:])
		if offset == 0 && dlen == -1 {
			return "", total_len
		} else if dlen == -1 {
			return strings.Join(s, ", "), total_len
		}

		next = offset + dlen
		s = append(s, v)
		offset = next
	}
	return strings.Join(s, ", "), total_len

}

func EscapeString(value []byte) string {
	valueStr := string(value)
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		valueStr = strings.Replace(valueStr, b, a, -1)
	}

	return valueStr
}
