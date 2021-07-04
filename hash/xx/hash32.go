package xx

import (
	"encoding/binary"
	"hash"
)

const (
	h32Prime1 = 0x9e3779b1
	h32Prime2 = 0x85ebca77
	h32Prime3 = 0xc2b2ae3d
	h32Prime4 = 0x27d4eb2f
	h32Prime5 = 0x165667b1
)

// Sum32 returns the xxHash32 checksum of the data.
func Sum32(p []byte) [4]byte {
	var h Hash32
	h.Reset()
	_, _ = h.Write(p)

	var sum [4]byte
	binary.BigEndian.PutUint32(sum[:], h.Sum32())
	return sum
}

// Hash32 is computes xxHash32 checksum.
type Hash32 struct {
	v1       uint32
	v2       uint32
	v3       uint32
	v4       uint32
	mem      [16]byte
	memSize  int
	totalLen uint32
	largeLen bool
}

var _ hash.Hash32 = (*Hash32)(nil)

// NewHash32 returns a new hash.Hash32 computing the xxHash32 checksum.
func NewHash32() *Hash32 {
	var h Hash32
	h.Reset()
	return &h
}

// Write is implementation of hash.Hash interface.
func (h *Hash32) Write(p []byte) (int, error) {
	pp := p
	h.totalLen += uint32(len(p))
	h.largeLen = h.largeLen || len(p) >= len(h.mem) || h.totalLen >= uint32(len(h.mem))

	// fill in temp buffer
	if h.memSize+len(p) < len(h.mem) {
		copy(h.mem[h.memSize:], p)
		h.memSize += len(p)
		return len(p), nil
	}

	// some data left from previous update
	if h.memSize > 0 {
		n := copy(h.mem[h.memSize:], p)
		h.v1 = h32Round(h.v1, readLE32(h.mem[0:4]))
		h.v2 = h32Round(h.v2, readLE32(h.mem[4:8]))
		h.v3 = h32Round(h.v3, readLE32(h.mem[8:12]))
		h.v4 = h32Round(h.v4, readLE32(h.mem[12:16]))
		p = p[n:]
		h.memSize = 0
	}

	for len(p) >= 16 {
		h.v1 = h32Round(h.v1, readLE32(p[0:4]))
		h.v2 = h32Round(h.v2, readLE32(p[4:8]))
		h.v3 = h32Round(h.v3, readLE32(p[8:12]))
		h.v4 = h32Round(h.v4, readLE32(p[12:16]))
		p = p[16:]
	}

	if len(p) > 0 {
		copy(h.mem[:], p)
		h.memSize = len(p)
	}

	return len(pp), nil
}

// Sum is implementation of hash.Hash interface.
func (h *Hash32) Sum(b []byte) []byte {
	var sum [4]byte
	binary.BigEndian.PutUint32(sum[:], h.Sum32())
	return append(b, sum[:]...)
}

// Reset is implementation of hash.Hash interface.
func (h *Hash32) Reset() {
	h.ResetSeed(0)
}

// Size is implementation of hash.Hash interface.
func (h *Hash32) Size() int {
	return 4
}

// BlockSize is implementation of hash.Hash interface.
func (h *Hash32) BlockSize() int {
	return 16
}

// Sum32 is implementation of hash.Hash32 interface.
func (h *Hash32) Sum32() uint32 {
	var sum uint32
	if h.largeLen {
		sum = rotl32(h.v1, 1) + rotl32(h.v2, 7) + rotl32(h.v3, 12) + rotl32(h.v4, 18)
	} else {
		sum = h.v3 + h32Prime5
	}
	sum += h.totalLen
	return h32Finalize(sum, h.mem[:h.memSize])
}

// ResetSeed is reset hash state with seed.
func (h *Hash32) ResetSeed(seed uint32) {
	*h = Hash32{
		v1: seed + h32Prime1 + h32Prime2,
		v2: seed + h32Prime2,
		v3: seed,
		v4: seed - h32Prime1,
	}
}

func h32Round(acc uint32, input uint32) uint32 {
	acc += input * h32Prime2
	acc = rotl32(acc, 13)
	acc *= h32Prime1
	return acc
}

func h32Finalize(sum uint32, p []byte) uint32 {
	for len(p) >= 4 {
		sum += readLE32(p) * h32Prime3
		sum = rotl32(sum, 17) * h32Prime4
		p = p[4:]
	}
	for len(p) > 0 {
		sum += uint32(p[0]) * h32Prime5
		sum = rotl32(sum, 11) * h32Prime1
		p = p[1:]
	}
	return h32Avalanche(sum)
}

func h32Avalanche(sum uint32) uint32 {
	sum ^= sum >> 15
	sum *= h32Prime2
	sum ^= sum >> 13
	sum *= h32Prime3
	sum ^= sum >> 16
	return sum
}
