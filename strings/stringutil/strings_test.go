package stringutil_test

import (
	"strings"
	"testing"

	xstrings "go.nanasi880.dev/x/strings/stringutil"
)

func TestHasPrefixFold(t *testing.T) {
	testSuites := []struct {
		s      string
		prefix string
		want   bool
	}{
		{
			s:      "Hello World",
			prefix: "hello",
			want:   true,
		},
		{
			s:      "Hello World",
			prefix: "Hello",
			want:   true,
		},
		{
			s:      "こんにちは世界",
			prefix: "こんにちは",
			want:   true,
		},
		{
			s:      "Hello World",
			prefix: "stub",
			want:   false,
		},
	}

	for i, suite := range testSuites {
		got := xstrings.HasPrefixFold(suite.s, suite.prefix)
		if got != suite.want {
			t.Fatal(i)
		}
	}
}

func TestHasSuffixFold(t *testing.T) {
	testSuites := []struct {
		s      string
		suffix string
		want   bool
	}{
		{
			s:      "Hello World",
			suffix: "world",
			want:   true,
		},
		{
			s:      "Hello World",
			suffix: "World",
			want:   true,
		},
		{
			s:      "こんにちは世界",
			suffix: "世界",
			want:   true,
		},
		{
			s:      "Hello World",
			suffix: "stub",
			want:   false,
		},
	}

	for i, suite := range testSuites {
		got := xstrings.HasSuffixFold(suite.s, suite.suffix)
		if got != suite.want {
			t.Fatal(i)
		}
	}
}

func BenchmarkHasPrefixFold(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 26; i++ {
		sb.WriteRune(rune('A' + i))
	}
	s := sb.String()

	b.Run("ToLower", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ss := strings.ToLower(s)
			strings.HasPrefix(ss, "a")
		}
	})
	b.Run("ToUpper", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ss := strings.ToUpper(s)
			strings.HasPrefix(ss, "A")
		}
	})
	b.Run("Fold", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			xstrings.HasPrefixFold(s, "a")
		}
	})
}
