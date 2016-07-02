package util

import (
	"errors"
)

const (
	MaxUInt64 = uint64(0xFFFFFFFFFFFFFFFF)
	MaxInt64 = int64(0x7FFFFFFFFFFFFFFF)
	MaxUInt64Digits = 20
	MaxInt64Digits = 19
	MaxUInt64OneLessOrderOfMag = MaxUInt64 /10
	MaxInt64OneLessOrderOfMag = MaxInt64 / 10
)

func BytesToInt64(b []byte) (int64, error) {
	neg := false
	if b[0] == '+' {
		b = b[1:]
	} else if b[0] == '-' {
		neg = true
		b = b[1:]
	}
	l := len(b)
	n := int64(0)
	if (l > MaxInt64Digits) {
		return 0, errors.New("value will overflow int64: " + string(b))
	} else if (l == MaxInt64Digits) {
		for i, v := range b {
			if v < '0' || v > '9' {
				return 0, errors.New("unknown value in byte: " + string(v) + " parsing " + string(b))
			}
			ld := int64(v - '0')
			if (i == 18) {
				//fmt.Printf("n=%v |  ld=%v\n", n, ld)
				if (n > MaxInt64OneLessOrderOfMag) {
					return 0, errors.New("value will overflow int64: " + string(b))
				} else if (n == MaxInt64OneLessOrderOfMag && ld > 7) {
					return 0, errors.New("value will overflow int64: " + string(b))
				}
			}
			n = n * 10 + ld
		}
	} else {
		// special optimized case - we overflow at 19 digits - anything less that than
		// is safe to relax the checking.
		for _, v := range b {
			if v < '0' || v > '9' {
				return 0, errors.New("unknown value in byte: " + string(v) + " parsing " + string(b))
			}
			n = n * 10 + int64(v - '0')
		}
	}

	if neg {
		return -n, nil
	}
	return n, nil
}

func BytesToUInt64(b []byte) (uint64, error) {
	if b[0] == '+' {
		b = b[1:]
	}
	l := len(b)
	n := uint64(0)
	if (l > MaxUInt64Digits) {
		return 0, errors.New("value will overflow uint64: " + string(b))
	} else if (l == MaxUInt64Digits) {
		for i, v := range b {
			if v < '0' || v > '9' {
				return 0, errors.New("unknown value in byte: " + string(v) + " parsing " + string(b))
			}
			ld := uint64(v - '0')
			if (i == 19) {
				//fmt.Printf("n=%v |  ld=%v\n", n, ld)
				if (n > MaxUInt64OneLessOrderOfMag) {
					return 0, errors.New("value will overflow uint64: " + string(b))
				} else if (n == MaxUInt64OneLessOrderOfMag && ld > 5) {
					return 0, errors.New("value will overflow uint64: " + string(b))
				}
			}
			n = n * 10 + ld
		}
	} else {
		// special optimized case - we overflow at 19 digits - anything less that than
		// is safe to relax the checking.
		for _, v := range b {
			if v < '0' || v > '9' {
				return 0, errors.New("unknown value in byte: " + string(v) + " parsing " + string(b))
			}
			n = n * 10 + uint64(v - '0')
		}
	}
	return n, nil
}