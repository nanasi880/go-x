package xx

import (
	"encoding/binary"
	"hash"
)

const (
	h64Prime1 uint64 = 0x9e3779b185ebca87
	h64Prime2 uint64 = 0xc2b2ae3d27d4eb4f
	h64Prime3 uint64 = 0x165667b19e3779f9
	h64Prime4 uint64 = 0x85ebca77c2b2ae63
	h64Prime5 uint64 = 0x27d4eb2f165667c5
)

// Sum64 returns the xxHash64 checksum of the data.
func Sum64(p []byte) [8]byte {
	var h Hash64
	h.Reset()
	_, _ = h.Write(p)

	var sum [8]byte
	binary.BigEndian.PutUint64(sum[:], h.Sum64())
	return sum
}

// Hash64 is computes xxHash64 checksum.
type Hash64 struct {
	v1       uint64
	v2       uint64
	v3       uint64
	v4       uint64
	mem      [32]byte
	memSize  int
	totalLen uint64
}

var _ hash.Hash64 = (*Hash64)(nil)

// NewHash64 returns a new hash.Hash64 computing the xxHash64 checksum.
func NewHash64() *Hash64 {
	var h Hash64
	h.Reset()
	return &h
}

// Write is implementation of hash.Hash interface.
func (h *Hash64) Write(p []byte) (int, error) {
	pp := p
	h.totalLen += uint64(len(p))

	// fill in temp buffer
	if h.memSize+len(p) < len(h.mem) {
		copy(h.mem[h.memSize:], p)
		h.memSize += len(p)
		return len(p), nil
	}

	// tmp buffer is full
	if h.memSize > 0 {
		n := copy(h.mem[h.memSize:], p)
		h.v1 = h64Round(h.v1, readLE64(h.mem[0:]))
		h.v2 = h64Round(h.v2, readLE64(h.mem[8:]))
		h.v3 = h64Round(h.v3, readLE64(h.mem[16:]))
		h.v4 = h64Round(h.v4, readLE64(h.mem[24:]))
		p = p[n:]
		h.memSize = 0
	}

	for len(p) >= 32 {
		h.v1 = h64Round(h.v1, readLE64(p[0:8]))
		h.v2 = h64Round(h.v2, readLE64(p[8:16]))
		h.v3 = h64Round(h.v3, readLE64(p[16:24]))
		h.v4 = h64Round(h.v4, readLE64(p[24:32]))
		p = p[32:]
	}

	if len(p) > 0 {
		copy(h.mem[:], p)
		h.memSize = len(p)
	}

	return len(pp), nil
}

// Sum is implementation of hash.Hash interface.
func (h *Hash64) Sum(b []byte) []byte {
	var sum [8]byte
	binary.BigEndian.PutUint64(sum[:], h.Sum64())
	return append(b, sum[:]...)
}

// Reset is implementation of hash.Hash interface.
func (h *Hash64) Reset() {
	h.ResetSeed(0)
}

// Size is implementation of hash.Hash interface.
func (h *Hash64) Size() int {
	return 8
}

// BlockSize is implementation of hash.Hash interface.
func (h *Hash64) BlockSize() int {
	return 32
}

// Sum64 is implementation of hash.Hash64 interface.
func (h *Hash64) Sum64() uint64 {
	var sum uint64
	if h.totalLen >= 32 {
		sum = rotl64(h.v1, 1) + rotl64(h.v2, 7) + rotl64(h.v3, 12) + rotl64(h.v4, 18)
		sum = h64MergeRound(sum, h.v1)
		sum = h64MergeRound(sum, h.v2)
		sum = h64MergeRound(sum, h.v3)
		sum = h64MergeRound(sum, h.v4)
	} else {
		sum = h.v3 + h64Prime5
	}
	sum += h.totalLen
	return h64Finalize(sum, h.mem[:h.memSize])
}

// ResetSeed is reset hash state with seed.
func (h *Hash64) ResetSeed(seed uint64) {
	*h = Hash64{
		v1: seed + h64Prime1 + h64Prime2,
		v2: seed + h64Prime2,
		v3: seed,
		v4: seed - h64Prime1,
	}
}

func h64Round(acc uint64, input uint64) uint64 {
	acc += input * h64Prime2
	acc = rotl64(acc, 31)
	acc *= h64Prime1
	return acc
}

func h64MergeRound(acc uint64, input uint64) uint64 {
	input = h64Round(0, input)
	acc ^= input
	acc = acc*h64Prime1 + h64Prime4
	return acc
}

func h64Finalize(sum uint64, p []byte) uint64 {
	for len(p) >= 8 {
		k1 := h64Round(0, readLE64(p))
		p = p[8:]
		sum ^= k1
		sum = rotl64(sum, 27)*h64Prime1 + h64Prime4
	}
	if len(p) >= 4 {
		sum ^= uint64(readLE32(p)) * h64Prime1
		p = p[4:]
		sum = rotl64(sum, 23)*h64Prime2 + h64Prime3
	}
	for len(p) > 0 {
		sum ^= uint64(p[0]) * h64Prime5
		p = p[1:]
		sum = rotl64(sum, 11) * h64Prime1
	}
	return h64Avalanche(sum)
}

func h64Avalanche(sum uint64) uint64 {
	sum ^= sum >> 33
	sum *= h64Prime2
	sum ^= sum >> 29
	sum *= h64Prime3
	sum ^= sum >> 32
	return sum
}
