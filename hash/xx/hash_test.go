package xx_test

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"testing"

	"go.nanasi880.dev/x/hash/xx"
	"go.nanasi880.dev/x/internal/testing/testutil"
)

func TestHash(t *testing.T) {
	testSuites := [...]struct {
		Value  []byte
		Sum32  []byte
		Sum64  []byte
		Sum128 []byte
	}{
		{
			Value:  []byte("Hello World !"),
			Sum32:  testutil.MustDecodeHexString(t, "ef104944"),
			Sum64:  testutil.MustDecodeHexString(t, "2f27a95e0277032b"),
			Sum128: testutil.MustDecodeHexString(t, "8c52e3056b8541c2780aae38ba5d77fa"),
		},
		{
			Value:  []byte("The quick brown fox jumps over the lazy dog"),
			Sum32:  testutil.MustDecodeHexString(t, "e85ea4de"),
			Sum64:  testutil.MustDecodeHexString(t, "0b242d361fda71bc"),
			Sum128: testutil.MustDecodeHexString(t, "ddd650205ca3e7fa24a1cc2e3a8a7651"),
		},
		{
			Value:  testutil.MustReadFile(t, "testdata/Square Polano.txt"),
			Sum32:  testutil.MustDecodeHexString(t, "64e3b0ab"),
			Sum64:  testutil.MustDecodeHexString(t, "2192c76c60a132f3"),
			Sum128: testutil.MustDecodeHexString(t, "eb22f44e32ac3f14c437688e07426857"),
		},
		{
			Value:  testutil.MustReadFile(t, "testdata/The Three-Cornered World.txt"),
			Sum32:  testutil.MustDecodeHexString(t, "703b5a25"),
			Sum64:  testutil.MustDecodeHexString(t, "624e25a34fe5e559"),
			Sum128: testutil.MustDecodeHexString(t, "9ca1941dfdfd1dd72f81241fcb240c15"),
		},
	}

	t.Run("Hash32", func(t *testing.T) {
		for suiteNo, suite := range testSuites {
			tag := fmt.Sprintf("suiteNo:%d", suiteNo)
			testHashMain(t, tag, xx.NewHash32(), suite.Value, suite.Sum32, func(value []byte) []byte {
				sum := xx.Sum32(value)
				return sum[:]
			})
		}
	})
	t.Run("Hash64", func(t *testing.T) {
		for suiteNo, suite := range testSuites {
			tag := fmt.Sprintf("suiteNo:%d", suiteNo)
			testHashMain(t, tag, xx.NewHash64(), suite.Value, suite.Sum64, func(value []byte) []byte {
				sum := xx.Sum64(value)
				return sum[:]
			})
		}
	})
	t.Run("Hash128", func(t *testing.T) {
		for suiteNo, suite := range testSuites {
			tag := fmt.Sprintf("suiteNo:%d", suiteNo)
			testHashMain(t, tag, xx.NewHash128(), suite.Value, suite.Sum128, func(value []byte) []byte {
				sum := xx.Sum128(value)
				return sum[:]
			})
		}
	})
}

func testHashMain(t *testing.T, tag string, h hash.Hash, value []byte, want []byte, sumFn func(value []byte) []byte) {
	var (
		hashSize  = h.Size()
		writeSize = h.BlockSize()/2 + 1
		sum       = sumFn(value)
	)
	if !bytes.Equal(sum, want) {
		testutil.Failf(t, "%s: sum:%s want:%s", tag, hex.EncodeToString(sum), hex.EncodeToString(want))
		return
	}
	for len(value) >= writeSize {
		n, err := h.Write(value[:writeSize])
		if err != nil {
			testutil.Failf(t, "%s: %v", tag, err)
			return
		}
		if n != writeSize {
			testutil.Failf(t, "%s: %v", tag, n)
			return
		}
		value = value[writeSize:]
	}
	if len(value) > 0 {
		n, err := h.Write(value)
		if err != nil {
			testutil.Failf(t, "%s: %v", tag, err)
			return
		}
		if n != len(value) {
			testutil.Failf(t, "%s: %v", tag, n)
			return
		}
	}
	sum = h.Sum(make([]byte, 0, hashSize))
	if len(sum) != hashSize {
		testutil.Failf(t, "%s: %v", tag, len(sum))
		return
	}
	if !bytes.Equal(sum, want) {
		testutil.Failf(t, "%s: sum:%s want:%s", tag, hex.EncodeToString(sum), hex.EncodeToString(want))
		return
	}
}

func BenchmarkHash(b *testing.B) {
	inputData := make([]byte, 255)
	for i := range inputData {
		inputData[i] = byte(i)
	}
	b.Run("xxHash32", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xx.Sum32(inputData)
		}
	})
	b.Run("xxHash64", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xx.Sum64(inputData)
		}
	})
	b.Run("xxHash128", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xx.Sum128(inputData)
		}
	})
	b.Run("CRC32", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = crc32.ChecksumIEEE(inputData)
		}
	})
	b.Run("CRC64", func(b *testing.B) {
		table := crc64.MakeTable(crc64.ISO)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = crc64.Checksum(inputData, table)
		}
	})
	b.Run("MD5", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = md5.Sum(inputData)
		}
	})
	b.Run("SHA1", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sha1.Sum(inputData)
		}
	})
	b.Run("SHA256", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = sha256.Sum256(inputData)
		}
	})
}
