package main

import (
	"fmt"
	"unsafe"
)

func f1() func() int {
	i := 12345

	return func() int {
		return i
	}
}

type funcvalStruct1 struct {
	fn uintptr
	i  int
}

func main() {
	f := f1()
	p := (**funcvalStruct1)(unsafe.Pointer(&f))
	fmt.Printf("%+v\n", *p)
	return
}
