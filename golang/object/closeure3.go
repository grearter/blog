package main

import (
	"fmt"
	"unsafe"
)

func f2() func() (int, int) {
	i := 12345
	j := 67890

	return func() (int, int) {
		return i, j
	}
}

type funcvalStruct2 struct {
	fn uintptr
	i  int
	j  int
}

func main() {
	f := f2()
	p := (**funcvalStruct2)(unsafe.Pointer(&f))
	fmt.Printf("%+v\n", *p)
	return
}
