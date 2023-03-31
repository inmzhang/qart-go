package qart_go

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"strconv"
	"strings"
)

func DataEncode(targetVersion int, ecl int, message string) (version int, codeWord [][]int) {

	mode := analyseMode(message)
	version = adjustVersion(targetVersion, mode, ecl, len(message))
	encodedMessage := ModeIndicator[mode] + charCountIndicator(version, mode, message) + encodeMethod[mode](message)

	encodedMessage = addTerminator(encodedMessage, version, ecl)
	encodedMessage = padMultipleOfEight(encodedMessage)
	encodedMessage = padToMaxCap(encodedMessage, version, ecl)

	dataCode := binaryStrToInt(encodedMessage)
	codeWord = groupDataCodewords(dataCode, version, ecl)
	return
}

var encodeMethod = [3]func(string) string{
	numericEncoding,
	alphaNumericEncoding,
	byteEncoding,
}

func analyseMode(message string) int {
	allInNumeric := true
	allInAlphanumeric := true
	for _, char := range message {
		if !strings.ContainsRune(NumericList, char) {
			allInNumeric = false
		}
		if !strings.ContainsRune(AlphanumericList, char) {
			allInAlphanumeric = false
		}
		if !allInNumeric && !allInAlphanumeric {
			return Byte
		}
	}
	if allInNumeric {
		return Numeric
	}
	return Alphanumeric
}

func adjustVersion(targetVersion int, mode int, ecl int, messageLength int) int {
	for version := targetVersion; version <= MaxVersion; version++ {
		if CharCap[ecl][version-1][mode] >= messageLength {
			return version
		}
	}
	return MaxVersion
}

func numericEncoding(message string) string {
	var encoded strings.Builder
	for i := 0; i < len(message); i += 3 {
		mLength := 3
		boundLength := len(message) - i
		if boundLength < 3 {
			mLength = boundLength
		}
		expectedNumBits := mLength*3 + 1
		intMessage, _ := strconv.Atoi(message[i : i+mLength])
		encoded.WriteString(fmt.Sprintf("%0*b", expectedNumBits, intMessage))
	}
	return encoded.String()
}

func alphaNumericEncoding(message string) string {
	var encoded strings.Builder
	var codeIndices []int
	for _, char := range message {
		codeIndices = append(codeIndices, strings.Index(AlphanumericList, string(char)))
	}
	numIndices := len(codeIndices)
	for i := 1; i < numIndices; i += 2 {
		intValue := codeIndices[i-1]*45 + codeIndices[i]
		encoded.WriteString(fmt.Sprintf("%011b", intValue))
	}
	if numIndices%2 != 0 {
		encoded.WriteString(fmt.Sprintf("%06b", codeIndices[numIndices-1]))
	}
	return encoded.String()
}

func byteEncoding(message string) string {
	encoded := strings.Builder{}
	latin1Encoder := charmap.ISO8859_1.NewEncoder()
	for _, char := range message {
		encodedChar, _ := latin1Encoder.Bytes([]byte{byte(char)})
		for _, b := range encodedChar {
			encoded.WriteString(fmt.Sprintf("%08b", b))
		}
	}
	return encoded.String()
}

func charCountIndicator(version int, mode int, message string) string {
	cciLengths := [][]int{
		{10, 9, 8, 8},
		{12, 11, 16, 10},
		{14, 13, 16, 12},
	}
	cciIndex := 2
	if version <= 9 {
		cciIndex = 0
	} else if version <= 26 {
		cciIndex = 1
	}
	return fmt.Sprintf("%0*b", cciLengths[cciIndex][mode], len(message))
}

func addTerminator(message string, version int, ecl int) string {
	maxBitsCap := 8 * BytesCap[version-1][ecl]
	numBitsPad := maxBitsCap - len(message)
	if numBitsPad < 4 {
		return message + strings.Repeat("0", numBitsPad)
	} else {
		return message + "0000"
	}
}

func padMultipleOfEight(message string) string {
	remainder := len(message) % 8
	if remainder != 0 {
		message += strings.Repeat("0", 8-remainder)
	}
	return message
}

func padToMaxCap(message string, version int, ecl int) string {
	maxBitsCap := 8 * BytesCap[version-1][ecl]
	for len(message) < maxBitsCap {
		if maxBitsCap-len(message) >= 16 {
			message += "1110110000010001"
		} else {
			message += "11101100"
		}
	}
	return message
}

func binaryStrToInt(message string) (code []int) {
	for i := 0; i < len(message); i += 8 {
		wordString := message[i : i+8]
		wordInt, _ := strconv.ParseInt(wordString, 2, 64)
		code = append(code, int(wordInt))
	}
	return
}

func groupDataCodewords(data []int, version int, ecl int) (codeWords [][]int) {
	characteristics := ErrorCorrectionCharacteristics[version-1][ecl]
	index := 0
	for i := 0; i < characteristics[0]; i++ {
		codeWords = append(codeWords, data[index:index+characteristics[1]])
		index += characteristics[1]
	}
	for i := 0; i < characteristics[2]; i++ {
		codeWords = append(codeWords, data[index:index+characteristics[3]])
		index += characteristics[3]
	}
	return
}
