package main

import (
	"fmt"
	"unsafe"
)

type funcvalStruct struct {
	fn uintptr
	// variable-size, fn-specific data here
}

func foo() {
	fmt.Println("func foo")
}

func main() {
	f := foo

	fmt.Printf("sizeof(f): %d\n", unsafe.Sizeof(f)) // sizeof(f): 8
	fmt.Printf("func foo addr: %p\n", foo)          // func foo addr: 0x1094ae0

	p := (**funcvalStruct)(unsafe.Pointer(&f))
	fmt.Printf("fn = 0x%x\n", (*p).fn) // fn = 0x1094ae0
	return
}
