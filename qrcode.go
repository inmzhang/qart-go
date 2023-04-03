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
	version, img := getRawQRCode(version, ecl, message)
	qrImg := img.(*image.NRGBA)
	if strings.HasSuffix(config.Path, "gif") {
		validateGif(savePath, config.Path)
		picImg := openGIF(config.Path)
		combinedGIF := ProcessGif(version, qrImg, picImg, config.Colorized, config.Contrast, config.Brightness)
		saveGIF(combinedGIF, savePath)

	} else if strings.HasSuffix(config.Path, "png") || strings.HasSuffix(config.Path, "jpeg") {
		picImg := openNRGBA(config.Path)
		CombineImages(version, qrImg, picImg, config.Colorized, config.Contrast, config.Brightness)
		saveImg := image.Image(qrImg)
		saveImg = ScaleImage(3, &saveImg)
		saveImage(&saveImg, savePath)
	}
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
		encodeErr = png.Encode(file, *img)
	} else {
		encodeErr = jpeg.Encode(file, *img, nil)
	}
	if encodeErr != nil {
		panic(encodeErr)
	}
}

func saveGIF(img *gif.GIF, savePath string) {
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

	encodeErr := gif.EncodeAll(file, img)
	if encodeErr != nil {
		panic(encodeErr)
	}
}

func validateGif(savePath, picPath string) {
	if strings.HasSuffix(picPath, "gif") && !strings.HasSuffix(savePath, "gif") {
		panic("You can only save a gif if you use a gif as the background!")
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

func openGIF(path string) *gif.GIF {
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

	img, err := gif.DecodeAll(file)
	if err != nil {
		panic(err)
	}
	return img
}
