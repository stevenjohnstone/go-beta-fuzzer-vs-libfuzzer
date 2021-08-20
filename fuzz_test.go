// +build gofuzzbeta

package fuzz

import "testing"

func FuzzBeta(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if magic(data) {
			t.Fatalf("magic is %v", data)
		}
	})
}

func FuzzLoopBeta(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if loopmagic(data) {
			t.Fatalf("magic is %v", data)
		}
	})
}
