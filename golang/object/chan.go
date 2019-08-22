package main

import (
	"fmt"
	"unsafe"
)

type chanStruct struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype unsafe.Pointer // element type
	sendx    uint           // send index
	recvx    uint           // receive index
}

func main() {
	c := make(chan int16, 10)
	c <- 1
	c <- 2
	close(c)

	fmt.Printf("sizeof(c): %d\n", unsafe.Sizeof(c))        // sizeof(c): 8
	fmt.Printf("len(c): %d, cap(c): %d\n", len(c), cap(c)) // len(c): 2, cap(c): 10

	p := (**chanStruct)(unsafe.Pointer(&c))
	fmt.Printf("p: %+v\n", *p) // p: &{qcount:2 dataqsiz:10 buf:0xc000094060 elemsize:2 closed:1 elemtype:0x10a5100 sendx:2 recvx:0}
	return
}
