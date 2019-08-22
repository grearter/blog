package main

import (
	"fmt"
	"unsafe"
)

func main() {
	m := map[string]string{
		"name": "golang",
	}

	fmt.Printf("m size: %d\n", unsafe.Sizeof(m))
	return
}
