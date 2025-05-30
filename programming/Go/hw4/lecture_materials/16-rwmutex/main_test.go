package main

import "testing"

func BenchmarkMutex(t *testing.B) {
	for i := 0; i < t.N; i++ {
		useMutex()
	}
}

func BenchmarkRWMutex(t *testing.B) {
	for i := 0; i < t.N; i++ {
		useRWMutex()
	}
}
