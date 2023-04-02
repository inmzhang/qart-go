package main

import (
	qart "github.com/inmzhang/qart-go"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var version int
	var ecl int
	var name string
	var dir string
	wd, _ := os.Getwd()

	app := &cli.App{
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
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() == 0 {
				log.Fatal("No message provided")
			}
			message := cCtx.Args().Get(0)
			qart.BasicQRCode(version, ecl, message, filepath.Join(dir, name))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
