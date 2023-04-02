package qart_go

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"math"
)

func CombineImages(version int, img1, img2 *image.NRGBA, colorized bool, contrast, brightness float64) {
	bgImg := imaging.AdjustContrast(img2, toPercentage(contrast))
	bgImg = imaging.AdjustBrightness(bgImg, toPercentage(brightness))
	bgImg = resizeAndFit(img1, bgImg)
	if !colorized {
		bgImg = deColorizeImage(bgImg)
	}
	mergeImages(version, img1, bgImg)
}

func ScaleImage(scale int, img *image.Image) (resized image.Image) {
	origWidth := (*img).Bounds().Dx()
	origHeight := (*img).Bounds().Dy()
	resized = imaging.Resize(*img, scale*origWidth, scale*origHeight, imaging.NearestNeighbor)
	return
}

func DrawMatrixPNG(matrix [][]int) image.Image {
	unitLen := 3
	x := 4 * unitLen
	y := 4 * unitLen
	matSize := len(matrix)
	pic := image.NewNRGBA(image.Rect(0, 0, (matSize+8)*unitLen, (matSize+8)*unitLen))
	for i := 0; i < pic.Bounds().Dx(); i++ {
		for j := 0; j < pic.Bounds().Dy(); j++ {
			pic.Set(i, j, image.White)
		}
	}

	for _, line := range matrix {
		for _, module := range line {
			if module != 0 {
				drawBlackUnit(pic, x, y, unitLen)
			}
			x += unitLen
		}
		x = 4 * unitLen
		y += unitLen
	}
	return pic
}

func drawBlackUnit(pic *image.NRGBA, x, y, unitLen int) {
	for i := 0; i < unitLen; i++ {
		for j := 0; j < unitLen; j++ {
			pic.Set(x+i, y+j, image.Black)
		}
	}
}

func toPercentage(val float64) float64 {
	return (val - 1) * 100
}

func resizeAndFit(img1, img2 *image.NRGBA) (resized *image.NRGBA) {
	width, height := img2.Bounds().Dx(), img2.Bounds().Dy()
	qrWidth, qrHeight := (*img1).Bounds().Dx(), (*img1).Bounds().Dy()
	ratio := float64(qrWidth) / float64(qrHeight)
	if width < height {
		resized = imaging.Resize(img2, qrWidth-24, (qrWidth-24)*int(1/ratio), imaging.Linear)
	} else {
		resized = imaging.Resize(img2, (qrHeight-24)*int(ratio), qrHeight-24, imaging.Linear)
	}
	return
}

func deColorizeImage(img *image.NRGBA) (decolorized *image.NRGBA) {
	decolorized = imaging.Grayscale(img)
	ditheringFS(decolorized)
	return
}

func mergeImages(version int, img1, img2 *image.NRGBA) {
	aligns := getAligns(version)
	qrWidth, qrHeight := img1.Bounds().Dx(), img1.Bounds().Dy()
	for i := 0; i < qrWidth-24; i++ {
		for j := 0; j < qrHeight-24; j++ {
			if i == 18 || i == 19 || i == 20 || j == 18 || j == 19 || j == 20 {
				continue
			}
			if i < 24 && j < 24 {
				continue
			}
			if i > qrWidth-49 && j < 24 {
				continue
			}
			if i < 24 && j > qrHeight-49 {
				continue
			}
			if inSlice(aligns, i, j) {
				continue
			}
			if i%3 == 1 && j%3 == 1 {
				continue
			}
			_, _, _, a := img2.At(i, j).RGBA()
			if a == 0 {
				continue
			}
			img1.Set(i+12, j+12, img2.At(i, j))
		}
	}
}

func getAligns(version int) (aligns [][2]int) {
	if version < 2 {
		return
	}
	aloc := AlignLocation[version-2]
	for a := 0; a < len(aloc); a++ {
		for b := 0; b < len(aloc); b++ {
			if (a == 0 && b == 0) || (a == 0 && b == len(aloc)-1) || (a == len(aloc)-1 && b == 0) {
				continue
			}
			for i := 3 * (aloc[a] - 2); i < 3*(aloc[a]+3); i++ {
				for j := 3 * (aloc[b] - 2); j < 3*(aloc[b]+3); j++ {
					aligns = append(aligns, [2]int{i, j})
				}
			}
		}
	}
	return
}

func inSlice(slice [][2]int, i, j int) bool {
	for _, v := range slice {
		if v[0] == i && v[1] == j {
			return true
		}
	}
	return false
}

func ditheringFS(img *image.NRGBA) {
	// Floyd-Steinberg dithering
	// https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering
	// Apply Floyd-Steinberg dithering
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the original pixel value
			oldColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			oldGray := float64(oldColor.Y)

			// Calculate the new pixel value
			newGray := 0
			if oldGray >= 128 {
				newGray = 255
			}

			// Set the new pixel value
			newColor := color.Gray{Y: uint8(newGray)}
			img.Set(x, y, newColor)

			// Calculate the error
			errorPixel := oldGray - float64(newGray)

			// Distribute the error to the neighboring pixels
			distributeError(img, x, y, errorPixel, 7.0/16.0)
			distributeError(img, x+1, y, errorPixel, 1.0/16.0)
			distributeError(img, x-1, y+1, errorPixel, 3.0/16.0)
			distributeError(img, x, y+1, errorPixel, 5.0/16.0)
		}
	}
}

func distributeError(img *image.NRGBA, x, y int, errorPixel, factor float64) {
	if x >= img.Bounds().Max.X || y >= img.Bounds().Max.Y {
		return
	}

	oldColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
	oldGray := float64(oldColor.Y)
	newGray := uint8(math.Max(0, math.Min(255, oldGray+errorPixel*factor)))
	newColor := color.Gray{Y: newGray}
	img.Set(x, y, newColor)
}
