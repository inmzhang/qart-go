package qart_go

import (
	"github.com/nfnt/resize"
	"image"
)

func ScaleImage(scale int, img *image.Image) (resized image.Image) {
	origWidth := (*img).Bounds().Dx()
	origHeight := (*img).Bounds().Dy()
	resized = resize.Resize(uint(scale*origWidth), uint(scale*origHeight), *img, resize.NearestNeighbor)
	return
}

func DrawMatrixPNG(matrix [][]int) image.Image {
	unitLen := 3
	x := 4 * unitLen
	y := 4 * unitLen
	matSize := len(matrix)
	pic := image.NewRGBA(image.Rect(0, 0, (matSize+8)*unitLen, (matSize+8)*unitLen))
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

func drawBlackUnit(pic *image.RGBA, x, y, unitLen int) {
	for i := 0; i < unitLen; i++ {
		for j := 0; j < unitLen; j++ {
			pic.Set(x+i, y+j, image.Black)
		}
	}
}
