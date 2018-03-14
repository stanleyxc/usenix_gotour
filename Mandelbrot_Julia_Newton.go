package main
//Mandelbrot using image/fractal
import (
	"fmt"
	"math/rand"						//for multi-color Mandelbrot
	"net/http"
	"stanleyxcnx/image"
	"stanleyxcnx/image/fractal"
)


func MandelbrotWeb(resp http.ResponseWriter, req *http.Request) {
	p := req.FormValue("p")				//X, Y, Co, Dc, iteration
	fmt.Printf("p: %v\n", p)
	var X, Y, i int
	var C0, Dc complex128

	n, err := fmt.Sscan(p, &X, &Y, &C0, &Dc, &i)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed.  Got %v: %v, %v, %v, %v, %v\n", n, X, Y, C0, Dc, i)


	colormap := image.ColorSet{In: make(map[int] image.RGBA, i), Out:make(map[int] image.RGBA, i)}
	colormap.In[0] = image.NewRGBA(255,255,255) 			// white for set.
	colormap.Out[0] = image.NewRGBA(0,0,0) 						// default black for outside of set.
	for j := 1; j < 100; j++ {
		//inColor := NewRGBA(j, j/2, j/4)
		//colormap.In[j] = inColor
		outColor := image.NewRGBA(rand.Intn(255),rand.Intn(255), rand.Intn(255) )			//random color
		colormap.Out[j] = outColor
	}
	mandelbrot := fractal.Mandelbrot(X,Y, C0, Dc, i, colormap)
	mandelbrot.Draw()
	fmt.Println("Mandelbrot image ready.")

	resp.Header().Set("Content-Type", "image/png")
	resp.Write(mandelbrot.Png())


}


func JuliaWeb(resp http.ResponseWriter, req *http.Request) {
	p := req.FormValue("p")					//x, y, Z0, Dz, C, iteration
	fmt.Printf("p: %v\n", p)
	var X, Y, i int
	var Z0, Dz, C complex128

	n, err := fmt.Sscan(p, &X, &Y, &Z0, &Dz, &C, &i)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed.  Got %v: %v, %v, %v, %v, %v, %v\n", n, X, Y, Z0, Dz, C, i)


	colormap := image.ColorSet{In: make(map[int] image.RGBA, i), Out:make(map[int] image.RGBA, i)}
	colormap.In[0] = image.NewRGBA(255,255,255) 			// white for set.
	colormap.Out[0] = image.NewRGBA(0,0,0) 						// default black for outside of set.
	for j := 1; j < 100; j++ {

		outColor := image.NewRGBA(rand.Intn(255),rand.Intn(255), rand.Intn(255) )			//random color
		colormap.Out[j] = outColor
	}
	julia := fractal.Julia(X,Y, Z0, Dz, C, i, colormap)
	julia.Draw()
	fmt.Println("Julia image ready.")

	resp.Header().Set("Content-Type", "image/png")
	resp.Write(julia.Png())

}

func NewtonWeb(resp http.ResponseWriter, req *http.Request) {
	p := req.FormValue("p")						// x, y, number, error margin, iteration, color (random or other.)
	fmt.Printf("p: %v\n", p)
	var X, Y, i, randomcolor int
	var v, e float64

	n, err := fmt.Sscan(p, &X, &Y, &v, &e, &i, &randomcolor)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed.  Got %v: %v, %v, %v, %v, %v, %v\n", n, X, Y, v, e, i, randomcolor)
	//return;
	colormap := make(map[int][3]image.RGBA, i)
	for j := 1; j < 100; j++ {
		if(randomcolor > 0) {
			c := image.NewRGBA(rand.Intn(255),rand.Intn(255), rand.Intn(255) )			//random color
			colormap[j] = [3]image.RGBA{c, c, c}									//same random color for each root.
		} else {
			red, green, blue := image.NewRGBA(255,0,0), image.NewRGBA(0,255,0), image.NewRGBA(0,0,255)
			red.R, red.A = uint8(255 - 255*j/i), uint8(255-j)						//scale color intensity proportional to iteration, fix alpha.
			green.G, green.A= uint8(255 - 255*j/i), uint8(255-j)
			blue.B, blue.A = uint8(255 - 255*j/i), uint8(255-j)
			colormap[j] = [3]image.RGBA{red, green, blue}								//colors: real root, complex, conjugate
		}
	}
	newton := fractal.Newton(X,Y, v, e, i, colormap)
	newton.Draw()
	fmt.Println("Newton Attractor image ready.")
	resp.Header().Set("Content-Type", "image/png")
	resp.Write(newton.Png())

}



func main() {

	fmt.Println("Fractal images on :4000")
	http.HandleFunc("/image/mandelbrot", MandelbrotWeb)
	http.HandleFunc("/image/julia", JuliaWeb)
	http.HandleFunc("/image/newton", NewtonWeb)
	http.ListenAndServe(":4000", nil)
}
