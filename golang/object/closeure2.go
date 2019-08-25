package main

import (
	"fmt"
	"unsafe"
)

func f3() func() (int, int) {
	i := 12345
	j := 67890

	return func() (int, int) {
		fmt.Printf("variable i addr: %p, variable j addr: %p\n", &i, &j)
		i++ // modify i
		j++ // modify j
		return i, j
	}
}

type funcvalStruct3 struct {
	fn uintptr
	i  uintptr
	j  uintptr
}

func main() {
	f := f3()
	fmt.Printf("f addr: 0x%p\n", f) // f addr: 0x0x1094a90

	f() // variable i addr: 0xc000094000, variable j addr: 0xc000094008

	p := (**funcvalStruct3)(unsafe.Pointer(&f))
	fmt.Printf("p.fn: 0x%x, p.i: 0x:%x, p.j: 0x%x\n", (*p).fn, (*p).i, (*p).j) // p.fn: 0x1094a90, p.i: 0x:c000094000, p.j: 0xc000094008
	return
}
