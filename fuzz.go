// +build !gofuzzbeta

package fuzz

func Fuzz(input []byte) int {
	if magic(input) {
		panic(input)
	}
	return 0
}

func FuzzLoop(input []byte) int {
	if loopmagic(input) {
		panic(input)
	}
	return 0
}
