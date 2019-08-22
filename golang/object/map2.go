package main

import (
	"fmt"
	"unsafe"
)

type mapStruct struct {
	count      int
	flags      uint8
	B          uint8
	noverflow  uint16
	hash0      uint32
	buckets    unsafe.Pointer
	oldbuckets unsafe.Pointer
	nevacuate  uintptr
	extra      unsafe.Pointer
}

func main() {
	m := make(map[string]string, 130)
	m["A"] = "AAA"
	m["B"] = "BBB"
	m["C"] = "CCC"
	m["D"] = "DDD"
	m["E"] = "EEE"

	fmt.Printf("sizeof(m): %d\n", unsafe.Sizeof(m))
	fmt.Printf("len(m): %d\n", len(m))

	p := (**mapStruct)(unsafe.Pointer(&m))

	fmt.Printf("p: %+v\n", *p)
}
