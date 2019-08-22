package main

import (
	"fmt"
	"unsafe"
)

type strStruct struct {
	str unsafe.Pointer
	len int
}

func main() {
	s := "hello"
	fmt.Printf("string size: %d\n", unsafe.Sizeof(s))

	p := (*strStruct)(unsafe.Pointer(&s))
	fmt.Printf("p.str: %p, p.len: %d\n", p.str, p.len)
}
