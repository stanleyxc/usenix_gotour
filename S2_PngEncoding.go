package main

import (
	"bytes"
	"net/http"
	"image"
	"image/color"
	"image/png"
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

/* copy from tour/pic ShowImage, removed base64 encoding. */
func (im Image) Png() []byte {
	var buf bytes.Buffer
	err := png.Encode(&buf, im)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (im Image) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "image/png")
	resp.Write(im.Png())
}

func main() {
	drawXY := func(x, y int) uint8 {
		return uint8(x*y)
	}
	drawCarrot := func(x, y int) uint8 {
		return uint8(x^y)
	}
	drawAvg := func(x, y int) uint8 {
		return uint8((x+y)/2)
	}

	a := Image{258,258, drawXY}
	b := Image{258,258, drawCarrot}
	c := Image{258,258, drawAvg}
	pic.ShowImage(a)

	http.Handle("/image/a", a)
	http.Handle("/image/b", b)
	http.Handle("/image/c", c)
	http.ListenAndServe(":4000", nil)


}
