package main

import (
	"fmt"
	"unsafe"
)

type sliceStruct struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func main() {
	words := make([]string, 0, 10)
	words = append(words, "hello", "golang", "world")

	fmt.Printf("slice size: %d\n", unsafe.Sizeof(words))

	p := (*sliceStruct)(unsafe.Pointer(&words))
	fmt.Printf("p.array: %p, p.len: %d, p.cap: %d\n", p.array, p.len, p.cap)
	return
}
