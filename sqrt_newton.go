//session 1 -  Nuts and Bolts

package main

import (
	"fmt"
	"strings"
	//"math"
)
//PrintSqrt prints square root of a number 10 times using Newton method - for code
func PrintSqrt(x float64) {
	z := 1.0
	for i:= 1; i <= 10; i++ {
		z -= (z*z -x) / (2*z)
		fmt.Printf("%d: %v\n", i, z)
	}
}
//PrintSqrtDz prints square root of a number x 10 times using Newton method with z initializer.
func PrintSqrtDz(x, z float64) {
	if z < 1.0 {
		z = 1.0
	}
	if z > 3.0 {
		z = 3.0
	}
	for i:= 1; i <= 10; i++ {
		z -= (z*z -x) / (2*z)
		fmt.Printf("%d: %v\n", i, z)
	}
}
//PrintSqrtDz prints square root of a number x, up to 10 times using Newton method with z initializer and t threshold value
// loop terminates once square root value achieves +/- t of x
func PrintSqrtDzt(x, z, t float64) {
	if z < 1.0 {
		z = 1.0
	}
	if z > 3.0 {
		z = 3.0
	}
	for i:= 1; i <= 10; i++ {
		z -= (z*z -x) / (2*z)
		fmt.Printf("%d: %v\n", i, z)
		if z1 := z*z - x; z*z > x && z1 <= t   {
			return
		}
		if z1 := x - z*z; z*z <= x && z1 <= t   {
			return
		}
	}
}
// PreciseSqrtDz find square root of x up to d decimal place precision using Newton method with initializer z
func PreciseSqrtDz(x, z float64, d int) float64 {
	if z < 1.0 {
		z = 1.0
	}
	if z > 3.0 {
		z = 3.0
	}
	//threshold := 1.0/math.Pow10(d)
	threshold := 1.0
	for _ = range make([]int, d) {
		threshold /= 10
	}
	for {
		z -= (z*z -x) / (2*z)
		if z1 := z*z - x; z*z > x && z1 <= threshold   {
			return z
		}
		if z1 := x - z*z; z*z <= x && z1 <= threshold   {
			return z
		}
	}
}
//Sqrt finds square root of a number using Newton method with precision upto d decimal place.
// eg. Sqrt(2, 6) -> 1.414213xxxx
func Sqrt(x float64, d int) float64 {
	z := 1.0
	//threshold := 1.0/math.Pow10(d)
	threshold := 1.0
	for _ = range make([]int, d) {
		threshold /= 10
	}
	//fmt.Println(threshold)
	for {
		z -= (z*z -x) / (2*z)
		if z1 := z*z - x; z*z > x && z1 <= threshold   {
			return z
		}
		if z1 := x - z*z; z*z <= x && z1 <= threshold   {
			return z
		}
	}
}

func main() {
	fmt.Printf("Square root of 2: %v (+/- 1e-06)\n", Sqrt(2, 6))
	PrintSqrt(2)
	PrintSqrtDz(2, 4)
	PrintSqrtDzt(2, 0.5, 0.000001)
	fmt.Printf("Square root of 2: %v (+/- 1e-06)\n", PreciseSqrtDz(2, 2.0, 6))

}
