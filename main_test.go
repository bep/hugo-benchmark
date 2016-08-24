package main

import (
	"testing"
)

func BenchmarkHugo(b *testing.B) {
	bench := createBench(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := bench.build(); err != nil {
			b.Fatal(err)
		}
	}
}
