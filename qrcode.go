package qart_go

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

type ArtConfig struct {
	Path       string
	Colorized  bool
	Contrast   float64
	Brightness float64
}

func ArtQRCode(version, ecl int, message, savePath string, config ArtConfig) int {
	validate(version, ecl, message)
	validateGif(savePath, config.Path)
	version, img := getRawQRCode(version, ecl, message)
	qrImg := img.(*image.NRGBA)
	picImg := openNRGBA(config.Path)
	CombineImages(version, qrImg, picImg, config.Colorized, config.Contrast, config.Brightness)
	saveImg := image.Image(qrImg)
	saveImg = ScaleImage(3, &saveImg)
	saveImage(&saveImg, savePath)
	return version
}

func BasicQRCode(version int, ecl int, message string, savePath string) int {
	validate(version, ecl, message)
	version, img := getRawQRCode(version, ecl, message)
	img = ScaleImage(3, &img)
	saveImage(&img, savePath)
	return version
}

func getRawQRCode(version int, ecl int, message string) (int, image.Image) {
	version, dataCodewords := DataEncode(version, ecl, message)
	ecc := ECCEncode(version, ecl, dataCodewords)
	structuredBits := StructureBits(version, ecl, dataCodewords, ecc)
	qrMatrix := QRMatrix(version, ecl, structuredBits)
	return version, DrawMatrixPNG(qrMatrix)
}

func saveImage(img *image.Image, savePath string) {
	file, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	err = file.Truncate(0)
	if err != nil {
		panic(err)
	}

	var encodeErr error
	if strings.HasSuffix(savePath, "png") {
		encodeErr = jpeg.Encode(file, *img, nil)
	} else if strings.HasSuffix(savePath, "jpeg") {
		encodeErr = png.Encode(file, *img)
	} else if strings.HasSuffix(savePath, "gif") {
		encodeErr = gif.Encode(file, *img, nil)
	} else {
		panic("Unsupported image format")
	}
	if encodeErr != nil {
		panic(encodeErr)
	}
}

func validateGif(savePath, picPath string) {
	if strings.HasSuffix(picPath, "gif") && !strings.HasSuffix(savePath, "gif") {
		panic("Unsupported image format")
	}
}

func validate(version int, ecl int, message string) {
	if !(version >= 1 && version <= 40) {
		panic("Wrong version! Please choose a value from 1 to 40!'")
	}
	if !(ecl >= 0 && ecl <= 3) {
		panic("Wrong Error-Correction-Level! Please choose a value from 0 to 3!'")
	}
	for _, char := range message {
		if !strings.ContainsRune(SupportedChars, char) {
			panic("Unsupported character: " + string(char))
		}
	}
}

func openNRGBA(path string) *image.NRGBA {
	if !strings.HasSuffix(path, "png") && !strings.HasSuffix(path, "jpeg") && !strings.HasSuffix(path, "gif") {
		panic("Unsupported image format")
	}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return img.(*image.NRGBA)
}
