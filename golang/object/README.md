# Go Object

## string
字符串是go中常用的基础数据类型之一，在golang源码src/runtime/string.go定义，如下：
```go
type stringStruct struct {
	str unsafe.Pointer
	len int
}
```
string本质上是一个struct, 包含2个成员变量
* str, 字符串的首地址
* len, 字符串的长度
<img src="https://github.com/grearter/blog/blob/master/golang/object/string.png" /><br/>
我们可以推断出, string变量的size为`8 + 8 = 16`

```go
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
	fmt.Printf("string size: %d\n", unsafe.Sizeof(s)) // string size: 16

	p := (*strStruct)(unsafe.Pointer(&s))
	fmt.Printf("p.str: %p, p.len: %d\n", p.str, p.len) // p.str: 0x10c7b14, p.len: 5
}
```

## slice
slice(数组切片)类型的底层同样是一个struct, 在源码中src/runtime/slice.go中定义, 如下:
```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```
<img src="https://github.com/grearter/blog/blob/master/golang/object/slice.png" /><br/>

slice struct包含3个成员变量:
* array, 指向数组切片内容的指针
* len, 切片存放的元素个数
* cap, 切片的总容量, cap >= len
我们可以推断出, string变量的size为`8 + 8 + 8 = 24`
```go
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

	fmt.Printf("slice size: %d\n", unsafe.Sizeof(words)) // slice size: 24

	p := (*sliceStruct)(unsafe.Pointer(&words))
	fmt.Printf("p.array: %p, p.len: %d, p.cap: %d\n", p.array, p.len, p.cap) // p.array: 0xc000096000, p.len: 3, p.cap: 10
	return
}
```

## map
map类型，在源码src/runtime/map.go中定义, 如下:
```go
// A header for a Go map.
type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}
```
hmap包含以下成员变量:
* count, 元素的个数
* flags, 状态标志
* B, 桶(bucket)容量 = 2 ^ B
* noverflow, 溢出buckets的个数
* hash0, 哈希种子
* buckets, 桶的地址
* oldBuckets, 旧桶的地址, 当map扩容时使用
* nevacuate, 搬迁进度, 小于nevacuate的已搬迁
* extra, 记录map的额外信息


先来看一下map变量的size:
```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	m := map[string]string{
		"name": "golang",
	}

	fmt.Printf("m size: %d\n", unsafe.Sizeof(m)) // m size: 8
	return
}
```
打印结果为`m size: 8`, 这是因为map变量实际是一个指针, 如下图:
<img src="https://github.com/grearter/blog/blob/master/golang/object/map.png" /><br/>
```go
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

	fmt.Printf("sizeof(m): %d\n", unsafe.Sizeof(m)) //  sizeof(m): 8
	fmt.Printf("len(m): %d\n", len(m)) // len(m): 5

	p := (**mapStruct)(unsafe.Pointer(&m)) 
	fmt.Printf("p: %+v\n", *p) // p: &{count:5 flags:0 B:5 noverflow:0 hash0:3540054924 buckets:0xc0000a4000 oldbuckets:<nil> nevacuate:0 extra:0xc0000b2000}
}
```

## chan
chan变量实际上也是`指针pointer`，其原型在src/runtime/chan.go中定义：
```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters
	lock mutex
}
```
chan struct包括以下成员变量:
* qcount, 循环队列(循环链表)中的元素个数
* dataqsiz, 循环队列的大小
* buf, 队列地址
* elemsize, 元素的size
* sendx, 循环队列的发送index
* recvx, 循环队列的接收index
* recvq, 等待接收的routine队列(双向链表)
* sendq, 等待发送的routine队列(双向链表)
* lock, 互斥锁
<img src="https://github.com/grearter/blog/blob/master/golang/object/chan.png" /><br/>
```go
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
```
详见: https://i6448038.github.io/2019/04/11/go-channel/

## function value
function value实际是`指针pointer`, src/runtime2.go:
```go
type funcval struct {
	fn uintptr
	// variable-size, fn-specific data here
}
```
<img src="https://github.com/grearter/blog/blob/master/golang/object/funcval.png" /><br/>

```go
package main

import (
	"fmt"
	"unsafe"
)

type funcvalStruct struct {
	fn uintptr
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
```
