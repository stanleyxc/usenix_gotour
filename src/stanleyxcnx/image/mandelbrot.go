package image

import (
  "math/cmplx"
)


type Mandelbrot struct {
	Pic *FancyImage
	Dx int
	Dy int
	C0 complex128
	Dc complex128
	Iteration int
	Colors ColorSet
}

func NewMandelbrot(x, y int, c0, dc complex128, it int) *Mandelbrot {
	im := NewFancyImage(x, y)
	return &Mandelbrot{im, x, y, c0, dc, it, ColorSet{}}

}
//Computes pixel value for a Mandelbrot plot.
func (m Mandelbrot) Compute(x, y int) RGBA {
		fx := float64(x)/float64(m.Dx)
		fy := float64(y)/float64(m.Dy)
		c := m.C0 + complex(real(m.Dc)*fx, imag(m.Dc)*fy)
		z := complex(0, 0)
		color := m.Colors.In[0]						//default set color
		for i:=0; i < m.Iteration; i++ {
			if z = z*z + c; cmplx.Abs(z) > 2 {
				color = m.Colors.Out[i]
				break
			}
		}
		return color
}
func (m Mandelbrot) DrawImage(cmap ColorSet) {
	m.Colors = cmap
  m.Pic.PixelDraw = m.Compute
	//m.Pic.SetDrawFun(m.Compute)
  m.Pic.Draw()
}
