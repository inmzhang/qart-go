package qart_go

import (
	"fmt"
	"strings"
)

func StructureBits(version int, ecl int, dataCodewords, ecc [][]int) (structuredBits string) {
	interleavedMessage := append(interleaveDC(version, ecl, dataCodewords), interleave(ecc)...)
	binaryStr := toBinaryStr(interleavedMessage)
	structuredBits = addRemainderBits(version, binaryStr)
	return
}

func interleaveDC(version int, ecl int, dc [][]int) (interleaved []int) {
	interleaved = interleave(dc)
	characteristic := ErrorCorrectionCharacteristics[version-1][ecl]
	if characteristic[3] != 0 {
		n := characteristic[2]
		for i := 0; i < n; i++ {
			block := dc[len(dc)+i-n]
			interleaved = append(interleaved, block[len(block)-1])
		}
	}
	return
}

func interleave(codewords [][]int) (interleaved []int) {
	var minLength int
	for _, block := range codewords {
		if len(block) < minLength || minLength == 0 {
			minLength = len(block)
		}
	}
	for i := 0; i < minLength; i++ {
		var unit []int
		for _, block := range codewords {
			unit = append(unit, block[i])
		}
		interleaved = append(interleaved, unit...)
	}
	return
}

func toBinaryStr(message []int) (str string) {
	var binaryStr []string
	for _, m := range message {
		binaryStr = append(binaryStr, fmt.Sprintf("%08b", m))
	}

	str = strings.Join(binaryStr, "")
	return
}

func addRemainderBits(version int, str string) string {
	numRemainderBits := NumRemainderBits[version-1]
	return str + strings.Repeat("0", numRemainderBits)
}
