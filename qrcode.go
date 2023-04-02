package qart_go

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func BasicQRCode(version int, ecl int, message string, savePath string) int {
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
