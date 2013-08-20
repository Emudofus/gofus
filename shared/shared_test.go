package shared

import (
	"testing"
	"math/rand"
	"time"
)

var (
	result string

	src = rand.NewSource(time.Now().UnixNano())
)

func BenchmarkNextString_8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rnd := NextString(src, 8)
		result = rnd
	}
}

func BenchmarkNextString_16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rnd := NextString(src, 16)
		result = rnd
	}
}

func BenchmarkNextString_32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rnd := NextString(src, 32)
		result = rnd
	}
}

func BenchmarkNextString_64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rnd := NextString(src, 64)
		result = rnd
	}
}

func BenchmarkNextString_512(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rnd := NextString(src, 512)
		result = rnd
	}
}

func BenchmarkNextString_4096(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rnd := NextString(src, 4096)
		result = rnd
	}
}
