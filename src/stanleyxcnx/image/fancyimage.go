package image

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	//"fmt"
)

//our own RGBA structure
type RGBA struct {
	R, G, B, A uint8
}
func NewRGBA(r, g, b int) RGBA {
	return RGBA{uint8(r), uint8(g), uint8(b), uint8(255)}		//fixed alpha
}
//a set of colors for Mandelbrot corresponding to iteration
type ColorSet struct {
	In map[int]  RGBA
	Out map[int] RGBA
}



type FancyImage struct {
	X int
	Y int
	Stride int
	Pixels []uint8 										// pixel data
	PixelDraw func(int, int) RGBA		// default function to generate pixel value.
}

func NewFancyImage(dx, dy int) *FancyImage {
	stride := 4*dx
	pixels := make([]uint8, 4*dx*dy)
	f := func(x, y int) RGBA {
		v := x*y
		return NewRGBA(v, v, v)
	}
	return &FancyImage{dx, dy, stride, pixels, f}
}



func (im FancyImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (im FancyImage) Width()  int {
	return im.X
}
func (im FancyImage) Height() int {
	return im.Y
}
func (im FancyImage) Bounds() image.Rectangle {
	return image.Rect(0,0, im.Width(), im.Height())
}
func (im FancyImage) At(x, y int) color.Color {

	if x <= im.Width() && y <= im.Height() {
		i := (y - im.Bounds().Min.Y)*im.Stride + (x-im.Bounds().Min.X)*4				//pixel offset for x, y coordinate.
		R, G, B, A := im.Pixels[i], im.Pixels[i+1],im.Pixels[i+2], im.Pixels[i+3]
		return color.RGBA{R, G, B, A}
	}
	return color.RGBA{0,0,0,255}
}


func (im FancyImage) Draw() {
	for y:=0; y < im.Y; y++ {
		for x:= 0; x < im.X; x++ {
			v := im.PixelDraw(x, y)					//pixel value: RGBA
			i := y*im.Stride + 4*x 				//pixel cell offset.
			im.Pixels[i] = v.R
			im.Pixels[i+1] = v.G
			im.Pixels[i+2] = v.B
			im.Pixels[i+3] = v.A
		}
	}
}

//no need for this.
func (im FancyImage) SetDrawFun(pixelDraw func(int, int) RGBA) {
	im.PixelDraw = pixelDraw
	//im.Draw()
	/*
	for y:=0; y < im.Y; y++ {
		for x:= 0; x < im.X; x++ {
			v := pixelDraw(x, y)					//pixel value: RGBA
			i := y*im.Stride + 4*x 				//pixel cell offset.
			im.Pixels[i] = v.R
			im.Pixels[i+1] = v.G
			im.Pixels[i+2] = v.B
			im.Pixels[i+3] = v.A
		}
	}
	*/
}


/* copy from tour/pic ShowImage, removed base64 encoding. */
func (im FancyImage) Png() []byte {
	var buf bytes.Buffer
	err := png.Encode(&buf, im)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (im FancyImage) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	im.Draw()
	resp.Header().Set("Content-Type", "image/png")
	resp.Write(im.Png())
}
