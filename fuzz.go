package fuzz

var leet = []byte{1, 3, 3, 7}

func magic(input []byte) bool {
	return len(input) == 4 && input[0] == leet[0] && input[1] == leet[1] && input[2] == leet[2] && input[3] == leet[3]
}

func loopmagic(input []byte) bool {
	if len(input) != 4 {
		return false
	}

	for i, v := range leet {
		if input[i] != v {
			return false
		}
	}
	return true
}

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
