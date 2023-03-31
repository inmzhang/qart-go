package qart_go

import (
	"fmt"
	"strconv"
	"strings"
)

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
