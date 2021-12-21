package asset

import (
	"embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed image map
var FS embed.FS

var ImgWhiteSquare = ebiten.NewImage(32, 32)

var ImgBrownBat = LoadImage("image/bat.png")

func init() {
	ImgWhiteSquare.Fill(color.White)
}

func LoadImage(p string) *ebiten.Image {
	f, err := FS.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	baseImg, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(baseImg)
}
