package xx

import (
	"encoding/binary"
	"fmt"
	"hash"
)

const (
	Hash128SecretSizeMin = 136
)
const (
	h128SecretDefaultSize   = 192
	h128StripeLen           = 64
	h128SecretConsumeLate   = 8
	h128InternalBufferSize  = 256
	h128SecretLastAccStart  = 7
	h128SecretMergeAccStart = 11
	h128MidSizeMax          = 240
	h128MidSizeStartOffset  = 3
	h128MidSizeLastOffset   = 17
)

var (
	h128DefaultSecret = [h128SecretDefaultSize]byte{
		0xb8, 0xfe, 0x6c, 0x39, 0x23, 0xa4, 0x4b, 0xbe, 0x7c, 0x01, 0x81, 0x2c, 0xf7, 0x21, 0xad, 0x1c,
		0xde, 0xd4, 0x6d, 0xe9, 0x83, 0x90, 0x97, 0xdb, 0x72, 0x40, 0xa4, 0xa4, 0xb7, 0xb3, 0x67, 0x1f,
		0xcb, 0x79, 0xe6, 0x4e, 0xcc, 0xc0, 0xe5, 0x78, 0x82, 0x5a, 0xd0, 0x7d, 0xcc, 0xff, 0x72, 0x21,
		0xb8, 0x08, 0x46, 0x74, 0xf7, 0x43, 0x24, 0x8e, 0xe0, 0x35, 0x90, 0xe6, 0x81, 0x3a, 0x26, 0x4c,
		0x3c, 0x28, 0x52, 0xbb, 0x91, 0xc3, 0x00, 0xcb, 0x88, 0xd0, 0x65, 0x8b, 0x1b, 0x53, 0x2e, 0xa3,
		0x71, 0x64, 0x48, 0x97, 0xa2, 0x0d, 0xf9, 0x4e, 0x38, 0x19, 0xef, 0x46, 0xa9, 0xde, 0xac, 0xd8,
		0xa8, 0xfa, 0x76, 0x3f, 0xe3, 0x9c, 0x34, 0x3f, 0xf9, 0xdc, 0xbb, 0xc7, 0xc7, 0x0b, 0x4f, 0x1d,
		0x8a, 0x51, 0xe0, 0x4b, 0xcd, 0xb4, 0x59, 0x31, 0xc8, 0x9f, 0x7e, 0xc9, 0xd9, 0x78, 0x73, 0x64,
		0xea, 0xc5, 0xac, 0x83, 0x34, 0xd3, 0xeb, 0xc3, 0xc5, 0x81, 0xa0, 0xff, 0xfa, 0x13, 0x63, 0xeb,
		0x17, 0x0d, 0xdd, 0x51, 0xb7, 0xf0, 0xda, 0x49, 0xd3, 0x16, 0x55, 0x26, 0x29, 0xd4, 0x68, 0x9e,
		0x2b, 0x16, 0xbe, 0x58, 0x7d, 0x47, 0xa1, 0xfc, 0x8f, 0xf8, 0xb8, 0xd1, 0x7a, 0xd0, 0x31, 0xce,
		0x45, 0xcb, 0x3a, 0x8f, 0x95, 0x16, 0x04, 0x28, 0xaf, 0xd7, 0xfb, 0xca, 0xbb, 0x4b, 0x40, 0x7e,
	}
	h128InitAcc = [8]uint64{
		h32Prime3,
		h64Prime1,
		h64Prime2,
		h64Prime3,
		h64Prime4,
		h32Prime2,
		h64Prime5,
		h32Prime1,
	}
)

// Sum128 returns the xxHash128 checksum of the data.
func Sum128(p []byte) [16]byte {
	var (
		sum [16]byte
		h   Hash128
	)
	h.Reset()
	_, _ = h.Write(p)
	digest := h.digest()
	binary.BigEndian.PutUint64(sum[:], digest.high)
	binary.BigEndian.PutUint64(sum[8:], digest.low)
	return sum
}

// Hash128 is computes xxHash128 checksum.
type Hash128 struct {
	acc             [8]uint64
	customSecret    [h128SecretDefaultSize]byte
	buffer          [h128InternalBufferSize]byte
	bufferedSize    int
	stripesSoFar    int
	totalLen        uint64
	stripesPerBlock int
	secretLimit     int
	seed            uint64
	extSecret       []byte
}

var _ hash.Hash = (*Hash128)(nil)

// NewHash128 returns a new Hash128 computing the xxHash128 checksum.
func NewHash128() *Hash128 {
	var h Hash128
	h.Reset()
	return &h
}

// Write is implementation of hash.Hash interface.
func (h *Hash128) Write(p []byte) (int, error) {
	pp := p
	secret := h.extSecret
	if len(secret) == 0 {
		secret = h.customSecret[:]
	}

	h.totalLen += uint64(len(p))

	// fill in tmp buffer
	if h.bufferedSize+len(p) <= len(h.buffer) {
		copy(h.buffer[h.bufferedSize:], p)
		h.bufferedSize += len(p)
		return len(p), nil
	}

	const stripes = h128InternalBufferSize / h128StripeLen

	if h.bufferedSize > 0 {
		n := copy(h.buffer[h.bufferedSize:], p)
		p = p[n:]
		h128ConsumeStripes(&h.acc, h.buffer[:], secret, h.secretLimit, &h.stripesSoFar, h.stripesPerBlock, stripes)
		h.bufferedSize = 0
	}

	if len(p) >= h128InternalBufferSize {
		for len(p)-h128InternalBufferSize > 0 {
			h128ConsumeStripes(&h.acc, p[:h128InternalBufferSize], secret, h.secretLimit, &h.stripesSoFar, h.stripesPerBlock, stripes)
			p = p[h128InternalBufferSize:]
		}
		// for last partial stripe
		offset := len(pp) - len(p)
		copy(h.buffer[len(h.buffer)-h128StripeLen:], pp[offset-h128StripeLen:])
	}

	n := copy(h.buffer[:], p)
	h.bufferedSize = n

	return len(pp), nil
}

// Sum is implementation of hash.Hash interface.
func (h *Hash128) Sum(b []byte) []byte {
	digest := h.digest()
	var sum [16]byte
	binary.BigEndian.PutUint64(sum[:], digest.high)
	binary.BigEndian.PutUint64(sum[8:], digest.low)
	return append(b, sum[:]...)
}

// Reset is implementation of hash.Hash interface.
func (h *Hash128) Reset() {
	secret := h128DefaultSecret[:]
	h.reset(0, secret, len(secret))
}

// Size is implementation of hash.Hash interface.
func (h *Hash128) Size() int {
	return 16
}

// BlockSize is implementation of hash.Hash interface.
func (h *Hash128) BlockSize() int {
	return 64
}

// Sum64 is implementation of hash.Hash64 interface.
func (h *Hash128) Sum64() uint64 {
	secret := h.extSecret
	if len(secret) == 0 {
		secret = h.customSecret[:]
	}
	if h.totalLen > h128MidSizeMax {
		digest := h.longDigest()
		return h128MergeAcc(&digest, secret[h128SecretMergeAccStart:], h.totalLen*h64Prime1)
	}
	if h.seed != 0 {
		return h128Sum64Seed(h.buffer[:h.bufferedSize], h.seed)
	}
	return h128Sum64Secret(h.buffer[:h.bufferedSize], secret[:h.secretLimit+h128StripeLen])
}

// ResetSeed is reset hash state with seed.
func (h *Hash128) ResetSeed(seed uint64) {
	if seed == 0 {
		h.Reset()
		return
	}
	if h.seed != seed {
		h128InitCustomSecret(h.customSecret[:], seed)
	}
	h.reset(seed, nil, h128SecretDefaultSize)
}

// ResetSecret is reset hash state with secret.
func (h *Hash128) ResetSecret(secret []byte) error {
	h.reset(0, secret, len(secret))
	if len(secret) < Hash128SecretSizeMin {
		return fmt.Errorf("secret too short")
	}
	return nil
}

func (h *Hash128) reset(seed uint64, secret []byte, secretSize int) {
	*h = Hash128{
		acc:             h128InitAcc,
		customSecret:    h.customSecret,
		buffer:          h.buffer,
		bufferedSize:    0,
		stripesSoFar:    0,
		totalLen:        0,
		stripesPerBlock: (secretSize - h128StripeLen) / h128SecretConsumeLate,
		secretLimit:     secretSize - h128StripeLen,
		seed:            seed,
		extSecret:       secret,
	}
}

func (h *Hash128) longDigest() [8]uint64 {
	secret := h.extSecret
	if len(secret) == 0 {
		secret = h.customSecret[:]
	}

	acc := h.acc
	if h.bufferedSize >= h128StripeLen {
		var (
			nbStripes      = (h.bufferedSize - 1) / h128StripeLen
			nbStripesSoFar = h.stripesSoFar
		)
		h128ConsumeStripes(&acc, h.buffer[:], secret, h.secretLimit, &nbStripesSoFar, h.stripesPerBlock, nbStripes)
		h128Accumulate512(&acc, h.buffer[h.bufferedSize-h128StripeLen:], secret[h.secretLimit-h128SecretLastAccStart:])
	} else {
		var lastStripe [h128StripeLen]byte
		catchupSize := h128StripeLen - h.bufferedSize
		copy(lastStripe[:], h.buffer[len(h.buffer)-catchupSize:])
		copy(lastStripe[catchupSize:], h.buffer[:h.bufferedSize])
		h128Accumulate512(&acc, lastStripe[:], secret[h.secretLimit-h128SecretLastAccStart:])
	}

	return acc
}

func (h *Hash128) digest() xxhUint128 {
	secret := h.extSecret
	if len(secret) == 0 {
		secret = h.customSecret[:]
	}

	if h.totalLen > h128MidSizeMax {
		var (
			acc  = h.longDigest()
			low  = h128MergeAcc(&acc, secret[h128SecretMergeAccStart:], h.totalLen*h64Prime1)
			high = h128MergeAcc(&acc, secret[h.secretLimit+h128StripeLen-64-h128SecretMergeAccStart:], ^(h.totalLen * h64Prime2))
		)
		return xxhUint128{high: high, low: low}
	}
	if h.seed != 0 {
		return h128Sum128Seed(h.buffer[:h.bufferedSize], h.seed)
	}
	return h128Sum128Secret(h.buffer[:h.bufferedSize], secret)
}

func h128ConsumeStripes(acc *[8]uint64, p []byte, secret []byte, secretLimit int, pStripesSoFar *int, stripesPerBlock int, stripes int) {
	stripesSoFar := *pStripesSoFar
	if stripesPerBlock-stripesSoFar <= stripes {
		var (
			stripesToEndOfBlock = stripesPerBlock - stripesSoFar
			stripesAfterBlock   = stripes - stripesToEndOfBlock
		)
		h128Accumulate(acc, p, secret[stripesSoFar*h128SecretConsumeLate:], stripesToEndOfBlock)
		h128ScrambleAcc(acc, secret[secretLimit:])
		h128Accumulate(acc, p[stripesToEndOfBlock*h128StripeLen:], secret, stripesAfterBlock)
		*pStripesSoFar = stripesAfterBlock
	} else {
		h128Accumulate(acc, p, secret[stripesSoFar*h128SecretConsumeLate:], stripes)
		*pStripesSoFar += stripes
	}
}

func h128InitCustomSecret(customSecret []byte, seed uint64) {
	_ = customSecret[h128SecretDefaultSize]

	const bounds = h128SecretDefaultSize / 16
	secret := h128DefaultSecret[:]
	for i := 0; i < bounds; i++ {
		low := readLE64(secret) + seed
		secret = secret[8:]
		writeLE64(customSecret, low)
		customSecret = customSecret[8:]

		high := readLE64(secret) - seed
		secret = secret[8:]
		writeLE64(customSecret, high)
		customSecret = customSecret[8:]
	}
}

func h128Accumulate(acc *[8]uint64, p []byte, secret []byte, nbStripes int) {
	for i := 0; i < nbStripes; i++ {
		h128Accumulate512(acc, p, secret)
		p = p[h128StripeLen:]
		secret = secret[h128SecretConsumeLate:]
	}
}

func h128Accumulate512(acc *[8]uint64, p []byte, secret []byte) {
	for i := range acc {
		var (
			val = readLE64(p)
			key = val ^ readLE64(secret)
		)
		p = p[8:]
		secret = secret[8:]

		acc[i^1] += val
		acc[i] += uint64(uint32(key&0xFFFFFFFF)) * uint64(uint32(key>>32))
	}
}

func h128ScrambleAcc(acc *[8]uint64, secret []byte) {
	for i := range acc {
		key := readLE64(secret[i*8:])
		val := acc[i]
		val = val ^ (val >> 47)
		val ^= key
		val *= h32Prime1
		acc[i] = val
	}
}

func h128MergeAcc(acc *[8]uint64, secret []byte, start uint64) uint64 {
	sum := start
	for i := 0; i < len(acc)/2; i++ {
		sum += h128MixAcc(acc[i*2:], secret[i*16:])
	}
	return h128Avalanche(sum)
}

func h128MixAcc(acc []uint64, secret []byte) uint64 {
	lhs := acc[0] ^ readLE64(secret)
	rhs := acc[1] ^ readLE64(secret[8:])
	return mul128Fold64(lhs, rhs)
}

func h128Avalanche(sum uint64) uint64 {
	sum = sum ^ (sum >> 37)
	sum *= 0x165667919e3779f9
	sum = sum ^ (sum >> 32)
	return sum
}

func h128rrmxmx(sum uint64, length uint64) uint64 {
	sum ^= rotl64(sum, 49) ^ rotl64(sum, 24)
	sum *= 0x9fb21c651e98df25
	sum ^= (sum >> 35) + length
	sum *= 0x9fb21c651e98df25
	return sum ^ (sum >> 28)
}

func h128Mix16B(p []byte, seed uint64, secret []byte) uint64 {
	var (
		low  = readLE64(p)
		high = readLE64(p[8:])
	)
	return mul128Fold64(
		(low^readLE64(secret))+seed,
		(high^readLE64(secret[8:]))-seed,
	)
}

func h128Mix32B(acc xxhUint128, input1 []byte, input2 []byte, secret []byte, seed uint64) xxhUint128 {
	acc.low += h128Mix16B(input1, seed, secret)
	acc.low ^= readLE64(input2) + readLE64(input2[8:])
	acc.high += h128Mix16B(input2, seed, secret[16:])
	acc.high ^= readLE64(input1) + readLE64(input1[8:])
	return acc
}

func h128Sum64Secret(p []byte, secret []byte) uint64 {
	switch {
	case len(p) <= 16:
		return h128Sum64Len0To16(p, 0, secret)
	case len(p) <= 128:
		return h128Sum64Len17To128(p, 0, secret)
	case len(p) <= h128MidSizeMax:
		return h128Sum64Len129To240(p, 0, secret)
	default:
		acc := h128InitAcc
		h128Sum64InternalLoop(&acc, p, secret)
		return h128MergeAcc(&acc, secret[h128SecretMergeAccStart:], uint64(len(p))*h64Prime1)
	}
}

func h128Sum64Seed(p []byte, seed uint64) uint64 {
	switch {
	case len(p) <= 16:
		return h128Sum64Len0To16(p, seed, h128DefaultSecret[:])
	case len(p) <= 128:
		return h128Sum64Len17To128(p, seed, h128DefaultSecret[:])
	case len(p) <= h128MidSizeMax:
		return h128Sum64Len129To240(p, seed, h128DefaultSecret[:])
	default:
		if seed == 0 {
			var (
				acc    = h128InitAcc
				secret = h128DefaultSecret[:]
			)
			h128Sum64InternalLoop(&acc, p, secret)
			return h128MergeAcc(&acc, secret[h128SecretMergeAccStart:], uint64(len(p))*h64Prime1)
		} else {
			var (
				acc    = h128InitAcc
				secret [h128SecretDefaultSize]byte
			)
			h128InitCustomSecret(secret[:], seed)
			h128Sum64InternalLoop(&acc, p, secret[:])
			return h128MergeAcc(&acc, secret[h128SecretMergeAccStart:], uint64(len(p))*h64Prime1)
		}
	}
}

func h128Sum64Len0To16(p []byte, seed uint64, secret []byte) uint64 {
	if len(p) > 8 {
		return h128Sum64Len9To16(p, seed, secret)
	}
	if len(p) >= 4 {
		return h128Sum64Len4To8(p, seed, secret)
	}
	if len(p) > 0 {
		return h128Sum64Len1To3(p, seed, secret)
	}
	return h64Avalanche(seed ^ readLE64(secret[56:]) ^ readLE64(secret[64:]))
}

func h128Sum64Len1To3(p []byte, seed uint64, secret []byte) uint64 {
	var (
		c1       = p[0]
		c2       = p[len(p)>>1]
		c3       = p[len(p)-1]
		combined = uint32(c1)<<16 | uint32(c2)<<24 | uint32(c3) | uint32(len(p))
		flip     = uint64(readLE32(secret)) ^ uint64(readLE32(secret[4:])) + seed
		keyed    = uint64(combined) ^ flip
	)
	return h64Avalanche(keyed)
}

func h128Sum64Len4To8(p []byte, seed uint64, secret []byte) uint64 {
	seed ^= uint64(swap32(uint32(seed))) << 32
	var (
		input1 = readLE32(p)
		input2 = readLE32(p[len(p)-4:])
		flip   = (readLE64(secret[8:]) ^ readLE64(secret[16:])) - seed
		input  = uint64(input2) + (uint64(input1) << 32)
		keyed  = input ^ flip
	)
	return h128rrmxmx(keyed, uint64(len(p)))
}

func h128Sum64Len9To16(p []byte, seed uint64, secret []byte) uint64 {
	var (
		flip1 = readLE64(secret[24:]) ^ readLE64(secret[32:]) + seed
		flip2 = readLE64(secret[40:]) ^ readLE64(secret[48:]) - seed
		low   = readLE64(p) ^ flip1
		high  = readLE64(p[len(p)-8:]) ^ flip2
		acc   = uint64(len(p)) + swap64(low) + high + mul128Fold64(low, high)
	)
	return h128Avalanche(acc)
}

func h128Sum64Len17To128(p []byte, seed uint64, secret []byte) uint64 {
	acc := uint64(len(p)) * h64Prime1
	if len(p) > 96 {
		acc += h128Mix16B(p[48:], seed, secret[96:])
		acc += h128Mix16B(p[len(p)-64:], seed, secret[112:])
	}
	if len(p) > 64 {
		acc += h128Mix16B(p[32:], seed, secret[64:])
		acc += h128Mix16B(p[len(p)-48:], seed, secret[80:])
	}
	if len(p) > 32 {
		acc += h128Mix16B(p[16:], seed, secret[32:])
		acc += h128Mix16B(p[len(p)-32:], seed, secret[48:])
	}
	acc += h128Mix16B(p, seed, secret)
	acc += h128Mix16B(p[len(p)-16:], seed, secret[16:])

	return h64Avalanche(acc)
}

func h128Sum64Len129To240(p []byte, seed uint64, secret []byte) uint64 {
	acc := uint64(len(p)) * h64Prime1
	for i := 0; i < 8; i++ {
		acc += h128Mix16B(p[i*16:], seed, secret[i*16:])
	}
	acc = h128Avalanche(acc)
	for i := 8; i < len(p)/16; i++ {
		acc += h128Mix16B(p[i*16:], seed, secret[((i-8)*16)+h128MidSizeStartOffset:])
	}
	acc += h128Mix16B(p[len(p)-16:], seed, secret[Hash128SecretSizeMin-h128MidSizeLastOffset:])
	return h128Avalanche(acc)
}

func h128Sum64InternalLoop(acc *[8]uint64, p []byte, secret []byte) {
	var (
		spritesPerBlock = (len(secret) - h128StripeLen) / h128SecretConsumeLate
		blockLen        = spritesPerBlock * h128StripeLen
		blocks          = (len(p) - 1) / blockLen
	)
	for i := 0; i < blocks; i++ {
		h128Accumulate(acc, p[i*blockLen:], secret, spritesPerBlock)
		h128ScrambleAcc(acc, secret[len(secret)-h128StripeLen:])
	}

	stripes := ((len(p) - 1) - (blockLen * blocks)) / h128StripeLen
	h128Accumulate(acc, p[blocks*blockLen:], secret, stripes)

	h128Accumulate512(acc, p[len(p)-h128StripeLen:], secret[len(secret)-h128StripeLen-h128SecretLastAccStart:])
}

func h128Sum128Secret(p []byte, secret []byte) xxhUint128 {
	switch {
	case len(p) <= 16:
		return h128Sum128Len0To16(p, 0, secret)
	case len(p) <= 128:
		return h128Sum128Len17To128(p, 0, secret)
	case len(p) <= h128MidSizeMax:
		return h128Sum128Len129To240(p, 0, secret)
	default:
		acc := h128InitAcc
		h128Sum64InternalLoop(&acc, p, secret)

		var h128 xxhUint128
		h128.low = h128MergeAcc(&acc, secret[h128SecretMergeAccStart:], uint64(len(p))*h64Prime1)
		h128.high = h128MergeAcc(&acc, secret[len(secret)-64-h128SecretMergeAccStart:], ^(uint64(len(p)))*h64Prime2)
		return h128
	}
}

func h128Sum128Seed(p []byte, seed uint64) xxhUint128 {
	switch {
	case len(p) <= 16:
		return h128Sum128Len0To16(p, seed, h128DefaultSecret[:])
	case len(p) <= 128:
		return h128Sum128Len17To128(p, seed, h128DefaultSecret[:])
	case len(p) <= h128MidSizeMax:
		return h128Sum128Len129To240(p, seed, h128DefaultSecret[:])
	default:
		if seed == 0 {
			acc := h128InitAcc
			secret := h128DefaultSecret[:]
			h128Sum64InternalLoop(&acc, p, secret)

			var h128 xxhUint128
			h128.low = h128MergeAcc(&acc, secret[h128SecretMergeAccStart:], uint64(len(p))*h64Prime1)
			h128.high = h128MergeAcc(&acc, secret[len(secret)-64-h128SecretMergeAccStart:], ^(uint64(len(p)))*h64Prime2)
			return h128
		} else {
			var (
				acc    = h128InitAcc
				secret [h128SecretDefaultSize]byte
			)
			h128InitCustomSecret(secret[:], seed)
			h128Sum64InternalLoop(&acc, p, secret[:])

			var h128 xxhUint128
			h128.low = h128MergeAcc(&acc, secret[h128SecretMergeAccStart:], uint64(len(p))*h64Prime1)
			h128.high = h128MergeAcc(&acc, secret[len(secret)-64-h128SecretMergeAccStart:], ^(uint64(len(p)))*h64Prime2)
			return h128
		}
	}
}

func h128Sum128Len0To16(p []byte, seed uint64, secret []byte) xxhUint128 {
	if len(p) > 8 {
		return h128Sum128Len9To16(p, seed, secret)
	}
	if len(p) >= 4 {
		return h128Sum128Len4To8(p, seed, secret)
	}
	if len(p) > 0 {
		return h128Sum128Len1To3(p, seed, secret)
	}
	var (
		flipLow  = readLE64(secret[64:]) ^ readLE64(secret[72:])
		flipHigh = readLE64(secret[80:]) ^ readLE64(secret[88:])
	)
	return xxhUint128{
		high: h64Avalanche(seed ^ flipHigh),
		low:  h64Avalanche(seed ^ flipLow),
	}
}

func h128Sum128Len1To3(p []byte, seed uint64, secret []byte) xxhUint128 {
	var (
		c1           = p[0]
		c2           = p[len(p)>>1]
		c3           = p[len(p)-1]
		combinedLow  = (uint32(c1) << 16) | (uint32(c2) << 24) | uint32(c3) | (uint32(len(p)) << 8)
		combinedHigh = rotl32(swap32(combinedLow), 13)
		flipLow      = uint64(readLE32(secret)^readLE32(secret[4:])) + seed
		flipHigh     = uint64(readLE32(secret[8:])^readLE32(secret[12:])) - seed
		keyedLow     = uint64(combinedLow) ^ flipLow
		keyedHigh    = uint64(combinedHigh) ^ flipHigh
	)
	var h128 xxhUint128
	h128.low = h64Avalanche(keyedLow)
	h128.high = h64Avalanche(keyedHigh)
	return h128
}

func h128Sum128Len4To8(p []byte, seed uint64, secret []byte) xxhUint128 {
	seed ^= uint64(swap32(uint32(seed))) << 32
	var (
		inputLow  = readLE32(p)
		inputHigh = readLE32(p[len(p)-4:])
		input     = uint64(inputLow) + (uint64(inputHigh) << 32)
		flip      = (readLE64(secret[16:]) ^ readLE64(secret[24:])) + seed
		keyed     = input ^ flip
	)
	m128 := mul64To128(keyed, h64Prime1+(uint64(len(p))<<2))
	m128.high += m128.low << 1
	m128.low ^= m128.high >> 3
	m128.low = m128.low ^ (m128.low >> 35)
	m128.low *= 0x9fb21c651e98df25
	m128.low = m128.low ^ (m128.low >> 28)
	m128.high = h128Avalanche(m128.high)
	return m128
}

func h128Sum128Len9To16(p []byte, seed uint64, secret []byte) xxhUint128 {
	var (
		flipLow   = (readLE64(secret[32:]) ^ readLE64(secret[40:])) - seed
		flipHigh  = (readLE64(secret[48:]) ^ readLE64(secret[56:])) + seed
		inputLow  = readLE64(p)
		inputHigh = readLE64(p[len(p)-8:])
		m128      = mul64To128(inputLow^inputHigh^flipLow, h64Prime1)
	)
	m128.low += (uint64(len(p)) - 1) << 54
	inputHigh ^= flipHigh
	m128.high += inputHigh + mul32To64(uint32(inputHigh), h32Prime2-1)
	m128.low ^= swap64(m128.high)

	h128 := mul64To128(m128.low, h64Prime2)
	h128.high += m128.high * h64Prime2
	h128.low = h128Avalanche(h128.low)
	h128.high = h128Avalanche(h128.high)
	return h128
}

func h128Sum128Len17To128(p []byte, seed uint64, secret []byte) xxhUint128 {
	var acc xxhUint128
	acc.low = uint64(len(p)) * h64Prime1

	if len(p) > 96 {
		acc = h128Mix32B(acc, p[48:], p[len(p)-64:], secret[96:], seed)
	}
	if len(p) > 64 {
		acc = h128Mix32B(acc, p[32:], p[len(p)-48:], secret[64:], seed)
	}
	if len(p) > 32 {
		acc = h128Mix32B(acc, p[16:], p[len(p)-32:], secret[32:], seed)
	}
	acc = h128Mix32B(acc, p, p[len(p)-16:], secret, seed)

	var h128 xxhUint128
	h128.low = acc.low + acc.high
	h128.high = (acc.low * h64Prime1) + (acc.high * h64Prime4) + ((uint64(len(p)) - seed) * h64Prime2)
	h128.low = h128Avalanche(h128.low)
	h128.high = -h128Avalanche(h128.high)
	return h128
}

func h128Sum128Len129To240(p []byte, seed uint64, secret []byte) xxhUint128 {
	var acc xxhUint128
	acc.low = uint64(len(p)) * h64Prime1

	for i := 0; i < 4; i++ {
		acc = h128Mix32B(acc, p[i*32:], p[i*32+16:], secret[i*32:], seed)
	}

	acc.low = h128Avalanche(acc.low)
	acc.high = h128Avalanche(acc.high)

	for i := 4; i < len(p)/32; i++ {
		acc = h128Mix32B(acc, p[i*32:], p[i*32+16:], secret[h128MidSizeStartOffset+(32*(i-4)):], seed)
	}

	// last bytes
	acc = h128Mix32B(acc, p[len(p)-16:], p[len(p)-32:], secret[Hash128SecretSizeMin-h128MidSizeLastOffset-16:], -seed)

	var h128 xxhUint128
	h128.low = acc.low + acc.high
	h128.high = (acc.low * h64Prime1) + (acc.high * h64Prime4) + ((uint64(len(p)) - seed) * h64Prime2)
	h128.low = h128Avalanche(h128.low)
	h128.high = -h128Avalanche(h128.high)
	return h128
}
