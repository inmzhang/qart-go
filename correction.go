package qart_go

func ECCEncode(version int, ecl int, dataCodeword [][]int) (ecc [][]int) {
	numECC := NumECCPerBlock[version-1][ecl]
	for _, dc := range dataCodeword {
		ecc = append(ecc, getECC(dc, numECC))
	}
	return
}

func getECC(dataCodeword []int, numECC int) (codeword []int) {
	generatorPoly := GeneratorPoly[numECC]
	codeword = dataCodeword
	for i := 0; i < len(dataCodeword); i++ {
		codeword = divide(codeword, generatorPoly)
	}
	return
}

func divide(messagePoly []int, generatorPoly []int) []int {
	if messagePoly[0] == 0 {
		return xor(messagePoly, make([]int, len(generatorPoly)))
	}
	for i := 0; i < len(generatorPoly); i++ {
		newVal := generatorPoly[i] + Log[messagePoly[0]]
		if newVal > 255 {
			newVal %= 255
		}
		generatorPoly[i] = PowerOfTwo[newVal]
	}
	return xor(messagePoly, generatorPoly)
}

func xor(messagePoly []int, generatorPoly []int) (remainder []int) {
	polyDiff := len(messagePoly) - len(generatorPoly)
	if polyDiff < 0 {
		messagePoly = append(messagePoly, make([]int, -polyDiff)...)
	} else if polyDiff > 0 {
		generatorPoly = append(generatorPoly, make([]int, polyDiff)...)
	}

	for i := 1; i < len(messagePoly); i += 1 {
		remainder = append(remainder, messagePoly[i]^generatorPoly[i])
	}
	return
}
