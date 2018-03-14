package main

import (
	"fmt"
	"strings"

)

func WordCount(s string) map[string]int {
	toks := strings.FieldsFunc(s, func(c rune) bool { return  c == ' ' || c == ',' || c== '(' || c== ')'  })
	f := make(map[string]int)
	for _, word := range toks {
		f[word]++					//default zero value of an int is 0
	/*
	//long version
		if r[word] > 0 {
			r[word]++
		} else {
			r[word] = 1
		}
	*/
	}
	return f
}
//return a list of words and location of repeated occurence.
func WordsAt(s string) map[string] []int {
	toks := strings.FieldsFunc(s, func(c rune) bool { return  c == ' ' || c == ',' || c== '(' || c== ')'  })
	f := make(map[string] []int)
	l := 0;
	for i, word := range toks {

		f[word] = append(f[word], l+i)
		l += len(word)
	}
	return f
}

//return a matrix of 8 bit unsigned integers
func Pic(dx, dy int) [][]uint8 {
	m := make([][]uint8, dy)
	for i, n := 0, dy; i < n; i++ {
		m[i] = make([]uint8, dx)
		for j, k :=0, dx; j < k; j++ {
			//m[i][j] = uint8(dx^dy)
			//m[i][j] = uint8((dx+dy)/2)
			m[i][j] = uint8(dx*dy)
		}

	}
	return m

}
//cube function

//high precision squart root
func main() {
	//input := "(id,created,employee(id,firstname,employeeType(id), d(kfdf), lastname),location), f(fdf)"
	input := "(id,created,employee(id,firstname,employeeType(id), lastname),location)"
	fmt.Printf("result: %v\n", WordCount(input))
	fmt.Printf("result: %v\n", WordsAt(input))
	input = " id created employee id  firstname employeeType id lastname location "
	fmt.Printf("result: %v\n", WordCount(input))
	fmt.Printf("result: %v\n", WordsAt(input))

	fmt.Printf("result: %v\n", Pic(4,5) )
}
