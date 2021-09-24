package fuzz

func FuzzLibFuzzer(input []byte) int {
	if magic(input) {
		panic(input)
	}
	return 0
}

func FuzzLoopLibFuzzer(input []byte) int {
	if loopmagic(input) {
		panic(input)
	}
	return 0
}
