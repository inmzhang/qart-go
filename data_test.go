package qart_go

import (
	"reflect"
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
			if got := numericEncoding(tt.message); got != tt.want {
				t.Errorf("numericEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_alphaNumericEncoding(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{name: "test 1", message: "AC-42", want: "0011100111011100111001000010"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alphaNumericEncoding(tt.message); got != tt.want {
				t.Errorf("alphaNumericEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_byteEncoding(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{name: "test 1", message: "adfasd12321AFSF&%", want: "0110000101100100011001100110000101110011011001000011000100110010001100110011001000110001010000010100011001010011010001100010011000100101"},
		{name: "test 2", message: "*&^*^&a?dSFSfasd123423/", want: "0010101000100110010111100010101001011110001001100110000100111111011001000101001101000110010100110110011001100001011100110110010000110001001100100011001100110100001100100011001100101111"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byteEncoding(tt.message); got != tt.want {
				t.Errorf("byteEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_charCountIndicator(t *testing.T) {
	type args struct {
		version int
		mode    int
		message string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test 1", args: args{10, 0, "1234567890"}, want: "000000001010"},
		{name: "test 2", args: args{3, 2, "1234567890$%WDAD"}, want: "00010000"},
		{name: "test 3", args: args{28, 1, "1234567890WDAD"}, want: "0000000001110"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := charCountIndicator(tt.args.version, tt.args.mode, tt.args.message); got != tt.want {
				t.Errorf("charCountIndicator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encode(t *testing.T) {
	type args struct {
		targetVersion int
		ecl           int
		message       string
	}
	tests := []struct {
		name         string
		args         args
		wantVersion  int
		wantCodeWord [][]int
	}{
		{
			"test 1",
			args{3, H, "adasdasda13123"},
			3,
			[][]int{
				{64, 230, 22, 70, 23, 54, 70, 23, 54, 70, 19, 19, 51},
				{19, 35, 48, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17},
			},
		},
		{
			"test 2",
			args{10, L, ">>::;;;13123?**@"},
			10,
			[][]int{
				{64, 1, 3, 227, 227, 163, 163, 179, 179, 179, 19, 51, 19, 35, 51, 242, 162, 164, 0, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236},
				{17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236},
				{17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17},
				{236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236, 17, 236},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion, gotCodeWord := DataEncode(tt.args.targetVersion, tt.args.ecl, tt.args.message)
			if gotVersion != tt.wantVersion {
				t.Errorf("DataEncode() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			}
			if !reflect.DeepEqual(gotCodeWord, tt.wantCodeWord) {
				t.Errorf("DataEncode() gotCodeWord = %v, want %v", gotCodeWord, tt.wantCodeWord)
			}
		})
	}
}
