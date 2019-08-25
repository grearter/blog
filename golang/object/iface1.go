package main

import (
	"fmt"
	"unsafe"
)

type Student struct {
	Id int
}

func (student Student) Print() {
	fmt.Println(student.Id)
	return
}

type Person interface {
	Print()
}

type ifaceStruct1 struct {
	tab  uintptr
	data unsafe.Pointer
}

func (i ifaceStruct1) Print() {
	panic("implement me")
}

func main() {
	var person Person = Student{Id: 12345}

	p := (*ifaceStruct1)(unsafe.Pointer(&person))

	fmt.Printf("p.tab: 0x:%x, p.data: 0x:%x\n", p.tab, p.data)

	data := (*Student)(p.data)
	fmt.Printf("data: %+v\n", *data)
}
