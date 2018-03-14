package fractal

import (
  "math/cmplx"
 "stanleyxcnx/image"
)

type Fractal struct {
	Name string
	*image.FancyImage          //embed
  MarginError float64
	Origin complex128					//Mandelbrot C0; Julia Z0
	Range complex128					//mandelbrot Dc; Julia Dz
	C complex128

	Iteration int
	Colors image.ColorSet
}

/*
func NewFractal(x, y int, c0, dc complex128, it int) *Fractal {
	im := NewFancyImage(x, y)
	c := complex(0,0)						// unused.
	return &Fractal{Name: "mandelbrot", im, Origin:c0, Range:dc, C:c,  Iteration:it, Colors:ColorSet{}}

}
*/


func Mandelbrot(x, y int, c0, dc complex128, it int, colors image.ColorSet) *Fractal {
  margin := 2.0
  mandelbrot := func(x0, y0 int) image.RGBA {
    fx := float64(x0)/float64(x)
		fy := float64(y0)/float64(y)
		c := c0 + complex(real(dc)*fx, imag(dc)*fy)
		z := complex(0, 0)
		color := colors.In[0]						//default set color
		for i:=0; i < it; i++ {
			if z = z*z + c; cmplx.Abs(z) > margin {
				color = colors.Out[i]
				break
			}
		}
		return color
  }
  im := image.NewFancyImage(x, y)
  im.PixelDraw = mandelbrot
	c := complex(0,0)
	return &Fractal{"Fractal.Mandelbrot",im, margin, c0, dc, c, it, colors}

}

func Julia(x, y int, z0, dz, C complex128, it int, colors image.ColorSet) *Fractal {
  margin := 2.0
  julia := func(x0, y0 int) image.RGBA {
    fx := float64(x0)/float64(x)
		fy := float64(y0)/float64(y)
		c := C
    z := z0 + complex(real(dz)*fx, imag(dz)*fy)
    color := colors.In[0]						//default set color
		for i:=0; i < it; i++ {
			if z = z*z + c; cmplx.Abs(z) > margin {
				color = colors.Out[i]
				break
			}
		}
		return color
  }
  im := image.NewFancyImage(x, y)
  im.PixelDraw = julia
	return &Fractal{"Fractal.Julia", im, margin, z0, dz, C, it, colors}
}


type NewtonAttractor struct {
  Name string
  *image.FancyImage
  Value, MarginError float64
  Iteration int
  Colors map[int] [3]image.RGBA		//colors based on iteration. [real, complex, conjugate]

}



//Newton finds cube root of real number only
func Newton(x, y int, n, e float64, it int, colors map[int] [3]image.RGBA) *NewtonAttractor {
  newton := func(x0, y0 int) image.RGBA {
    fx := float64(x0) - float64(x)/2										//shift coordinate so the center of plot is the middle of the generated image.
    fy := float64(y0) - float64(y)/2
    sx := float64(x)/2*n																//scale factor
    sy := float64(y)/2*n																//scale factor
    z := complex(fx/sx, fy/sy)
    v := complex(n, 0)

    var i = 0;
    for  {
      if i >= it {
        return image.NewRGBA(255, 255, 255)				//return white pixel if exceeded max iteration
      }
      z -= (z*z*z - v) / (3*z*z)
      delta := cmplx.Abs(z*z*z  - v);
      if delta <= e {
        break
      }
      i++
    }

    if re, img := real(z), imag(z); re < 0 {			//complex root: the real part of any complext roots is a  negative number
    //if re, img := real(z), imag(z); math.Abs(re*re*re - n) >= e {			//complex root
      if img < 0 {    //conjugate
        return colors[i][2]
      }
      return colors[i][1]
    }
    return colors[i][0]
  }

  im := image.NewFancyImage(x, y)
  im.PixelDraw = newton
	return &NewtonAttractor{"Fractal.NewtonAttractor",im, n, e, it, colors}

}
