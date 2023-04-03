package mpc

import (
	"testing"
)

func TestSplitAndCombain(t *testing.T) {
	secrets := make([]byte, 256)
	for i := 0; int(i) < 256; i++ {
		secrets[i] = byte(i)
	}

	parts := SplitSecret(secrets, 10, 4)
	selected := parts[4:8]
	recovered := CombainParts(selected)

	for i := 0; i < len(secrets); i++ {
		if secrets[i] != recovered[i] {
			t.Fail()
		}
	}
}

func BenchmarkSplit(b *testing.B) {
	secrets := make([]byte, 256)
	for i := 0; i < 256; i++ {
		secrets[i] = byte(i)
	}

	for i := 0; i < b.N; i++ {
		SplitSecret(secrets, 10, 5)
	}
}

func BenchmarkCombain(b *testing.B) {
	secrets := make([]byte, 256)
	for i := 0; i < 256; i++ {
		secrets[i] = byte(i)
	}

	parts := SplitSecret(secrets, 10, 4)
	selected := parts[2:6]

	for i := 0; i < b.N; i++ {
		CombainParts(selected)
	}
}
