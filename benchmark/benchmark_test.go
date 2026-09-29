package benchmark

import (
	"testing"
)

func BenchmarkPut(b *testing.B) {
	store :=newMemtable()
	
}