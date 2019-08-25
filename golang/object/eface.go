package main

import (
	"fmt"
	"unsafe"
)

type efaceStruct struct {
	_type uintptr
	data  unsafe.Pointer
}

func main() {
	n := 12345
	fmt.Printf("variable n addr: %p\n", &n)

	var i interface{} = n
	p := (*efaceStruct)(unsafe.Pointer(&i))
	data := *((*int)(p.data))
	fmt.Printf("p._type: 0x%x, p.data: 0x%x\n", p._type, p.data)
	fmt.Printf("data: %d\n", data)
	return
}
