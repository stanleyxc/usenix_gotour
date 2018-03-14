//session 3 - Equivalent Binary Trees
package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

func Walk(t *tree.Tree, ch chan int)  {
	ch <- t.Value
	if t.Left != nil {
		 Walk(t.Left, ch)
	}
	if t.Right != nil {
	   Walk(t.Right, ch)
	}
}
func Same(t1, t2 *tree.Tree) bool {
	c1, c2 := make(chan int), make(chan int)
	m := make(map[int]int)						//use this map as comparator
	go func () {
		Walk(t1, c1)
		close(c1)					//close channel after walk done.
	}()
	go func () {
		Walk(t2, c2)
		close(c2)
	}()

	//read from first channel
	for v, next:= <-c1; next; v, next = <-c1 {
		m[v]++
  }
	//read from second channel.
	for v, next:= <-c2; next; v, next = <-c2 {
		m[v]--
  }
	//every element in m should have a value of  zero, otherwise trees aren't the same.
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true;
}

func main() {
	c := make(chan int)
	go func () {
		Walk(tree.New(1), c)
		close(c)
	}()

	for v, next := <-c; next; v, next = <-c {
		fmt.Printf("Tree Element: %v\n", v)
	}
	fmt.Println("Done printing tree")

	t1, t2 := tree.New(1), tree.New(1)
	fmt.Printf("Trees are the same?: %v\n", Same(t1, t2))
	t3, t4 := tree.New(1), tree.New(2)
	fmt.Printf("Trees are the same?: %v\n", Same(t3, t4))

}
