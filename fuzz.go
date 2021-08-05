package fuzz

func magic(input []byte) bool {
	return len(input) == 4 && input[0] == 1 && input[1] == 3 && input[2] == 3 && input[3] == 7
}

func Fuzz(input []byte) int {
	if magic(input) {
		panic("found")
	}
	return 0
}
