package main

import (
	//"io"
	//"strings"

	//"encoding/base64"
	"fmt"
	//"math"
	"math/rand"						//for multi-color Mandelbrot
	"math/cmplx"
	"bytes"
	"net/http"

	"image"
	"image/color"
	"image/png"
	"golang.org/x/tour/pic"
)


type RGBA struct {
	R, G, B, A uint8
}
func NewRGBA(r, g, b int) RGBA {
	return RGBA{uint8(r), uint8(g), uint8(b), uint8(255)}		//fix alpha
}
//color map for colorful Mandelbrot
type ColorMap struct {
	In map[int]RGBA
	Out map[int]RGBA
}



type Image struct {
	X int
	Y int
	Stride int
	Pixels []uint8 										// pixel data
	PixelDraw func(int, int) RGBA		// default function to generate pixel value.
}
//first constructor
func NewImage(dx, dy int) *Image {
	stride := 4*dx
	pixels := make([]uint8, 4*dx*dy)
	f := func(x, y int) RGBA {
		v := x*y
		return NewRGBA(v, v, v)
	}
	return &Image{dx, dy, stride, pixels, f}
}

func (im Image) Draw(pixelDraw func(int, int) RGBA) {
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

	if x <= im.Width() && y <= im.Height() {
		i := (y - im.Bounds().Min.Y)*im.Stride + (x-im.Bounds().Min.X)*4				//pixel offset for x, y coordinate.
		R, G, B, A := im.Pixels[i], im.Pixels[i+1],im.Pixels[i+2], im.Pixels[i+3]
		return color.RGBA{R, G, B, A}
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


type Mandelbrot struct {
	Pic *Image
	Dx int
	Dy int
	C0 complex128
	Dc complex128
	Iteration int
	Colors ColorMap
}
func NewMandelbrot(x, y int, c0, dc complex128, it int) *Mandelbrot {
	im := NewImage(x, y)
	return &Mandelbrot{im, x, y, c0, dc, it, ColorMap{}}

}
func (m Mandelbrot) Compute(x, y int) RGBA {
		fx := float64(x)/float64(m.Dx)
		fy := float64(y)/float64(m.Dy)
		c := m.C0 + complex(real(m.Dc)*fx, imag(m.Dc)*fy)
		z := complex(0, 0)
		color := m.Colors.In[0]
		for i:=0; i < m.Iteration; i++ {
			if z = z*z + c; cmplx.Abs(z) > 2 {

				color = m.Colors.Out[i]
				break
			}
		}
		return color
}
func (m Mandelbrot) Draw(cmap ColorMap) {
	m.Colors = cmap
	m.Pic.Draw(m.Compute)
}


func main() {
	drawXY := func(x, y int) RGBA {
		v := x*y
		return NewRGBA(v, v, v)
	}
	drawCarrot := func(x, y int) RGBA {
		v := x^y
		return NewRGBA(v, v, v)
	}
	drawAvg := func(x, y int) RGBA {
		v:= (x+y)/2
		return NewRGBA(v, v, v)
	}

	a := NewImage(258,258)
	a.Draw(drawXY)
	b := NewImage(258,258)
	b.Draw(drawCarrot)
	c := NewImage(258,258)
	c.Draw(drawAvg)
	man := NewMandelbrot(1024,1024, complex(-2.5, -2.5), complex(5, 5), 100)
	colormap := ColorMap{In: make(map[int]RGBA, 100), Out:make(map[int]RGBA, 100)}

	colormap.In[0] = NewRGBA(255,255,255) 			// white for set.
	colormap.Out[0] = NewRGBA(0,0,0) 						// default black for outside of set.
	//multi-color for edges
	for j := 1; j < 100; j++ {
		//inColor := NewRGBA(j, j/2, j/4)
		//colormap.In[j] = inColor
		outColor := NewRGBA(rand.Intn(255),rand.Intn(255), rand.Intn(255) )			//random color
		colormap.Out[j] = outColor
	}

	man.Draw(colormap)

	pic.ShowImage(a)
	pic.ShowImage(man.Pic)
	fmt.Println("Serving images on :4000")
	http.Handle("/image/a", a)
	http.Handle("/image/b", b)
	http.Handle("/image/c", c)
	http.Handle("/image/man", man.Pic)
	http.ListenAndServe(":4000", nil)




}
