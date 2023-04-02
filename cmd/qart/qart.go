package main

import qart "github.com/inmzhang/qart-go"

func main() {
	version := 20
	ecl := qart.H
	message := "https://zhuanlan.zhihu.com/p/387753099"
	qart.BasicQRCode(version, ecl, message, "qrcode.png")
}
