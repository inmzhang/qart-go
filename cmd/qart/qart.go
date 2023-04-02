package main

import (
	qart "github.com/inmzhang/qart-go"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var version, ecl int
	var name, dir, picPath string
	var colorized bool
	var contrast, brightness float64

	wd, _ := os.Getwd()

	app := &cli.App{
		Name:        "qart",
		Description: "A go implementation of QR-Code generator and artistic picture embedding.",
		Suggest:     true,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "version",
				Value:       1,
				Aliases:     []string{"v"},
				Usage:       "The version means the length of a side of the QR-Code picture. From little size to large is 1 to 40.",
				Destination: &version,
			},
			&cli.IntFlag{
				Name:        "level",
				Value:       qart.H,
				Aliases:     []string{"l"},
				Usage:       "Use this argument to choose an Error-Correction-Level: 0(Low), 1(Medium) or 2(Quartile), 3(High).",
				Destination: &ecl,
			},
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Usage:       "The filename of output tailed with one of {'.jpg', '.png', '.bmp', '.gif'}.",
				Value:       "qrcode.png",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "dir",
				Aliases:     []string{"d"},
				Usage:       "The directory of output, default to current directory.",
				Value:       wd,
				Destination: &dir,
			},
			&cli.StringFlag{
				Name:        "picture",
				Aliases:     []string{"p"},
				Usage:       "The path to the picture to be embedded in the QR-Code.",
				Destination: &picPath,
			},
			&cli.BoolFlag{
				Name:        "colorized",
				Aliases:     []string{"c"},
				Usage:       "Whether to colorize the QR-Code.",
				Destination: &colorized,
			},
			&cli.Float64Flag{
				Name:        "contrast",
				Aliases:     []string{"con"},
				Value:       1.0,
				Usage:       "A floating point value controlling the enhancement of contrast. Factor 1.0 always returns a copy of the original image, lower factors mean less color (brightness, contrast, etc), and higher values more.",
				Destination: &contrast,
			},
			&cli.Float64Flag{
				Name:        "brightness",
				Aliases:     []string{"bri"},
				Value:       1.0,
				Usage:       "A floating point value controlling the enhancement of brightness. Factor 1.0 always returns a copy of the original image, lower factors mean less color (brightness, contrast, etc), and higher values more.",
				Destination: &brightness,
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() == 0 {
				log.Fatal("No message provided, please input a message to be encoded.")
			}
			message := cCtx.Args().Get(0)
			if picPath != "" {
				config := qart.ArtConfig{Path: picPath, Colorized: colorized, Contrast: contrast, Brightness: brightness}
				qart.ArtQRCode(version, ecl, message, filepath.Join(dir, name), config)
				return nil
			}
			qart.BasicQRCode(version, ecl, message, filepath.Join(dir, name))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	//config := qart.ArtConfig{Path: "C:\\GoProgramming\\qart-go\\cmd\\qart\\surface.png", Colorized: true, Contrast: 1.5, Brightness: 1.0}
	//qart.ArtQRCode(1, 3, "https://zhuanlan.zhihu.com/p/387753099", "qrcode.png", config)
}
