package qart_go

import (
	"fmt"
	"strconv"
	"strings"
)

func analyseMode(message string) (mode int) {
	allInNumeric := true
	allInAlphanumeric := true
	for _, char := range message {
		if !strings.ContainsRune(NumericList, char) {
			allInNumeric = false
		}
		if !strings.ContainsRune(AlphanumericList, char) {
			allInAlphanumeric = false
		}
	}
	if allInNumeric {
		mode = Numeric
	} else if allInAlphanumeric {
		mode = Alphanumeric
	} else {
		mode = Byte
	}
	return
}

func adjustVersion(version int, mode int, ecl int, messageLength int) int {
	for i := 0; i < 40; i++ {
		if CharCap[ecl][i][mode] > messageLength && version < i+1 {
			return i + 1
		}
	}
	return 40
}

func numericalEncoding(message string) (string, error) {
	var encoded strings.Builder
	for i := 0; i < len(message); i += 3 {
		mLength := 3
		boundLength := len(message) - i
		if boundLength < 3 {
			mLength = boundLength
		}
		expectedNumBits := mLength*3 + 1
		intMessage, err := strconv.Atoi(message[i : i+mLength])
		if err != nil {
			return "", fmt.Errorf("failed to convert message slice to integer: %w", err)
		}
		encoded.WriteString(fmt.Sprintf("%0*b", expectedNumBits, intMessage))
	}
	return encoded.String(), nil
}
