package main

import (
	"fmt"
	"testing"
)

func BenchmarkSprintf(b *testing.B) {
	key := "key"
	value := "value"
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%s:%s", key, value)
	}
}

func BenchmarkConcat(b *testing.B) {
	key := "key"
	value := "value"
	for i := 0; i < b.N; i++ {
		_ = key + ":" + value
	}
}
