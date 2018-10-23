package main

import "fmt"

func mutit(j []int) {
	for i := range j {
		j[i] = 44
	}
}
func main() {
	j := make([]int, 4)
	fmt.Println("j = ", j)
	mutit(j)
	fmt.Println("j = ", j)
}
