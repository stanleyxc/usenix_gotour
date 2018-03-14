package image

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"fmt"
)

type Image struct {
	X int
	Y int
	PixelDraw func(int, int) uint8		// default function to generate pixel value.
}

func NewImage(x, y int, draw func(int, int) uint8) *Image {
  return &Image{x, y, draw}
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
// Simple At function just compute Pixel value.
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
		fmt.Printf("%v\n", err)
		panic(err)
	}
	//fmt.Println("image encoded")
	b := buf.Bytes()
	//fmt.Printf("image data lenth: %v at  %v, %v\n", len(b), im.X, im.Y)
	return b
}


func (im Image) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "image/png")
	resp.Write(im.Png())
	fmt.Printf("image (%v, %v, %v) served\n", im.X, im.Y)
}

//usage
// var im Image
// var b []bytes
// im = NewImage(1, 1, func(x, y int) uint8 { return uint8((x+y)/2) })
// b = im.Png()
