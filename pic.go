//session 1 -  Nuts and Bolts

package main

import (
	"fmt"
	"strings"
	//"math"
)
////Slices
//Pic generates dx X dy matrix of uint8 values
func Pic(dx, dy int) [][]uint8 {
	matrix := make([][]uint8, dy)
	for i, n := 0, dy; i < n; i++ {
		matrix[i] = make([]uint8, dx)
		for j, k :=0, dx; j < k; j++ {
			//m[i][j] = uint8(dx^dy)
			//m[i][j] = uint8((dx+dy)/2)
			matrix[i][j] = uint8(dx*dy)
		}
	}
	return matrix
}

//PicFun generates dx X dy matrix of uint8 values using f function
func PicFun(dx, dy int, f func(int, int) uint8) [][]uint8 {
	matrix := make([][]uint8, dy)
	for i, n := 0, dy; i < n; i++ {
		matrix[i] = make([]uint8, dx)
		for j, k :=0, dx; j < k; j++ {
			matrix[i][j] = f(dx, dy)
		}
	}
	return matrix
}


func main() {

	fmt.Printf("Pic Matrix: %v\n", Pic(3, 5))
	var fxy = func(x, y int) uint8 {
		return uint8(x*y)
	}
	var favg = func(x, y int) uint8 {
		return uint8((x+y)/2)
	}
	fmt.Printf("Pic Matrix (x*y): %v\n", PicFun(3, 5, fxy))
	fmt.Printf("Pic Matrix ((x+y)/2): %v\n", PicFun(3, 5, favg))
	fmt.Printf("Pic Matrix (x^y): %v\n", PicFun(3, 5, func(x, y int) uint8 { return uint8(x^y) } ) )
}
