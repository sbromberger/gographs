package main

import "fmt"

type Foo struct {
	a []uint32
	x int
}

func newFoo(n int) Foo {
	a := make([]uint32, 0, n)
	return Foo{a, n}
}

func doFoo(f Foo) {
	f.a = append(f.a, 20)
	f.x = 20
}

func doFooP(f *Foo) {
	f.a = append(f.a, 30)
	f.x = 30
}

func main() {
	f := newFoo(10)
	f.a = append(f.a, 10)
	fmt.Println("f = ", f)
	doFoo(f)
	fmt.Println("f = ", f)
	doFooP(&f)
	fmt.Println("f = ", f)
}
