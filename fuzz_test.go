// +build gofuzzbeta

package fuzz

import "testing"

func Fuzz(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if magic(data) {
			t.Fatalf("magic is %v", data)
		}
	})
}

func FuzzLoop(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if loopmagic(data) {
			t.Fatalf("magic is %v", data)
		}
	})
}
