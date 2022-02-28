package byteutil

import (
	"fmt"
)

const (
	KB  Size = 1000
	MB       = 1000 * KB
	GB       = 1000 * MB
	TB       = 1000 * GB
	PB       = 1000 * TB
	EB       = 1000 * PB
	KiB Size = 1024
	MiB      = 1024 * KiB
	GiB      = 1024 * MiB
	TiB      = 1024 * GiB
	PiB      = 1024 * TiB
	EiB      = 1024 * PiB
)

var (
	sizeToStringMap = map[Size]string{
		KB:  "KB",
		MB:  "MB",
		GB:  "GB",
		TB:  "TB",
		PB:  "PB",
		EB:  "EB",
		KiB: "KiB",
		MiB: "MiB",
		GiB: "GiB",
		TiB: "TiB",
		PiB: "PiB",
		EiB: "EiB",
	}
)

type Size int64

func (i Size) Int64() int64 {
	return int64(i)
}

func (i Size) Format(state fmt.State, verb rune) {

	if verb != 'v' && verb != 's' {
		// not support
		return
	}

	switch {
	case i == 0:
		_, _ = fmt.Fprintf(state, "0B")
	case i < 0:
		factors := [...]Size{EB, PB, TB, GB, MB, KB}
		if state.Flag('+') {
			factors = [...]Size{EiB, PiB, TiB, GiB, MiB, KiB}
		}
		var factor Size
		for _, f := range factors {
			if -f >= i {
				factor = f
				break
			}
		}
		if factor != 0 {
			_, _ = fmt.Fprintf(state, "%.02f%s", float64(i)/float64(factor), sizeToStringMap[factor])
		} else {
			_, _ = fmt.Fprintf(state, "%dB", i.Int64())
		}
	case i > 0:
		factors := [...]Size{KB, MB, GB, TB, PB, EB}
		if state.Flag('+') {
			factors = [...]Size{KiB, MiB, GiB, TiB, PiB, EiB}
		}

		if i < factors[0] {
			_, _ = fmt.Fprintf(state, "%dB", i.Int64())
			return
		}

		var factor Size
		for index := 1; index < len(factors); index++ {
			if i < factors[index] {
				factor = factors[index-1]
				break
			}
		}
		if factor == 0 {
			factor = factors[len(factors)-1]
		}
		_, _ = fmt.Fprintf(state, "%.02f%s", float64(i)/float64(factor), sizeToStringMap[factor])
	}
}
