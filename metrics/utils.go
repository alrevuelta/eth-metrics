package metrics

import (
	"strconv"
)

func BoolToUint64(in bool) uint64 {
	if in {
		return uint64(1)
	}
	return uint64(0)
}

func ToBytes48(x []byte) [48]byte {
	var y [48]byte
	copy(y[:], x)
	return y
}

func UToStr(x uint64) string {
	return strconv.FormatUint(x, 10)
}
