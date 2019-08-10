# Go(1.12) GC

## 进程的内存布局

一个进程的虚拟内存由多个`段(segment)`组成，如下图所示:
<img src="https://github.com/grearter/blog/blob/master/golang/gc/memory_layout.png" />

* 内核空间(Kernel): 内核总是驻留在内存中，并映射到进程的虚拟内存中，但不允许进程读写操作
* 参数列表(argv)与环境变量(environ)
* 栈(Stack): 主要存储函数的参数与局部变量等，由操作系统自动分配与释放
* 堆(heap): 用于运行时动态分配内存
* 为初始化的全局变量(bss)
* 初始化的全局变量(initialized data)
* 代码段(text)

## 什么是GC
GC(Garbage Collection), 即垃圾回收，是一种自动内存管理机制，负责回收不再被进程使用的对象占用的内存。

## 手动管理内存 vs 自动管理内存

#### 手动管理内存
在C程序中，开发者显式的调用`malloc`函数来动态分配内存，在使用完成之后，必须显式显式调用`free`函数来释放动态申请的内存；
在C++程序中，显式调用`new`操作符来动态申请内存，显式调用`delete`来释放动态申请的内存；

手动内存管理需要开发者时刻注意对象的生命周期:
  * 显式调用malloc或delete来释放内存
  * 当一块内存释放之后，可能还需要清空指向已释放内存的指针，以避免出现`野指针`造成程序崩溃
  * 不能过早的回收还在使用的内存
  * 调用第三方lib时，要明确对象的所有权(明确对象由谁来进行释放, 避免重复释放)
  * 容易造成内存泄漏
 
### 自动管理内存
通过GC可实现自动管理内存
* GC可以解决大部分的内存泄漏问题
* GC对只会对未被引用的对象进行回收，不会出现野指针与重复释放
* 无需关心第三方lib对象的问题, 这降低来与其它模块的耦合
* 开发人员只需专注于业务

## GC vs 资源回收
GC是一种内存管理机制，回收的是不再被引用的对象的内存。
对于进程的其它资源如socket、file等，GC是不负责回收的，需要开发者在程序中显式的调用相关函数来进行回收。

## Go 1.12 GC原理
### 术语说明
#### SWT
Stop The World，在GC的一些阶段中，需要暂停所有用户goroutine, 来确保所有的P达到GC安全点(GC safe-point)

#### Root对象
Root对象指不需要其它对象就能直接访问到的对象。主要包括`栈变量`、`全局变量`以及其它`堆外(off-heap)内存变量`。

#### 可达性
即通过Root对象可以`直接`或`间接`访问到。
<img src="https://github.com/grearter/blog/blob/master/golang/gc/reachable_objs.png" />
一般来说，如果一个对象是`不可达`的, 那么此对象是需要被GC回收。

#### 标记和清扫(Mark and Sweep)
标记: 将`可达`对象进行标记 <br/>
清扫: 将`不可达`对象进行回收

#### Span
TODO

### GC流程
go 1.12.1 使用`写屏障的并发标记和清除`来进行垃圾回收。

#### 1. 执行清除终止(sweep termination)
    a. Stop The World, 确保所有P达到GC安全点
    b. 清除任何未清除过的span。只有在预期时间之前强制执行此GC周期时，才会有未清除的span
#### 2. 执行标记阶段(mark)
    a. 为标记做准备: 将gcphase设置为`_GCmark`, 启用`写屏障(write barrier)`, 启用`mutator assist`, 并将`Root对象`排入队列(enqueue)
       在所有的P都启用`写屏障`之前, 不会扫描任何对象(这是STW完成的)
    b. Start The World, 从现在开始，GC工作`标记worker(工作由调度器启动的)`和assists(作为allocation的一部分)来完成。
       写屏障将覆写的指针和任何指针写的新指针值都着色。在这之后, `写屏障`会把`重写指针(overwrite pointer)`与`新的指针值(new pointer value)`进行着色。
       新分配的对象被立即标记的`black`
    c. GC执行Root对象标记工作。扫描所有`栈变量`, 着色所有全局变量，以及着色堆外运行时数据结构中的任何堆指针。
       扫描(scan)`栈(stack)`对暂停对应的goroutine, 扫描完成之后回复goroutine
    d. 耗尽`worker队列`中的灰色对象, 扫描每个'gray对象':
        i. 将`gray对象`标记为`black`
        ii. 对在该对象中找到的所有指针进行着色并放入`worker队列`中
    e. 由于 GC work 分散在本地缓存中，因此 GC 使用分布式终止算法来检测何时不再有根标记作业或灰色对象（参见gcMarkDone函数）。此时, GC 状态转换到标记终止(mark termination)。
#### 3. 标记终止(mark termination)
    a. Stop The World
    b. 将gcphase设置为`_GCmarktermination`, 并禁用 workers 和 assists
    c. 进行内务整理，如flushing mcaches(runtime/mcache.go)
    
#### 4. 清扫阶段(sweep)
    a. 为清扫阶段做准备:
        i. 将gcphase设置为`_GCoff`
        ii. 设置清除状态
        iii. 禁用`写屏障(write barrier)`
    b. Start The World, 从现在开始, 所有新申请的对象都记为`white`, 如有必要，在使用spans前清除spans 
    c. GC 在后台进行并发清除并响应allocation

### 并发清除(concurrent sweep)
清除阶段与正常程序执行并发进行。<br/>
在后台`sweep goroutine`中，堆内存(heap)被惰性(当goroutine需要另一个span时)且并发地逐个span扫描(这有助于非CPU密集型的程序) <br/>
在STW标记终止的结尾，所有的span都被标记为需要清除。<br/>
后台`sweep goroutine`简单地逐个清除span。
