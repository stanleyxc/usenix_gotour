package main

import (
	"image"
	"image/color"
	"golang.org/x/tour/pic"
)

//session 2

type Image struct {
	X int
	Y int
	PixelDraw func(int, int) uint8		// default function to generate pixel value.
}

func (im Image) Width()  int {
	return im.X
}
func (im Image) Height() int {
	return im.Y
}
func (im Image) Bounds() image.Rectangle {
	return image.Rect(0,0, im.Width(), im.Height())
}
//simple At function that just compute pixel value user-suppied draw function.
func (im Image) At(x, y int) color.Color {
	if x <= im.X && y <= im.Y {
		v := im.PixelDraw(x, y);
		return color.RGBA{v, v, v, 255}
	}
	return color.RGBA{0,0,0,255}
}

func (im Image) ColorModel() color.Model {
	return color.RGBAModel
}

func main() {
	drawf := func(x, y int) uint8 {
		return uint8(x*y)
	}
	g := Image{277,277, drawf}
	pic.ShowImage(g)


}
