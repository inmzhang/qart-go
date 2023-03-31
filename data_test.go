package qart_go

import (
	"testing"
)

func Test_analyseMode(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		wantMode int
	}{
		{name: "test numeric mode", message: "1233141", wantMode: Numeric},
		{name: "test alphanumeric mode", message: "1233141ABJK$", wantMode: Alphanumeric},
		{name: "test byte mode", message: "1233141ABJK?!asdasd%^#$%^&#", wantMode: Byte},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMode := analyseMode(tt.message); gotMode != tt.wantMode {
				t.Errorf("analyseMode() = %v, want %v", gotMode, tt.wantMode)
			}
		})
	}
}

func Test_numericalEncoding(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{name: "test 1", message: "01234567", want: "000000110001010110011000011"},
		{name: "test 2", message: "0123456789012345", want: "000000110001010110011010100110111000010100111010100101"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := numericalEncoding(tt.message); got != tt.want {
				t.Errorf("numericalEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}
