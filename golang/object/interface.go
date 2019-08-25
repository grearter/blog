package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var i interface{}
	i = 123

	fmt.Printf("sizeo(i): %d\n", unsafe.Sizeof(i))
}
