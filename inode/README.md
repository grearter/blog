# inode

## 什么是inode

在Linux系统中，每个文件(包括目录)都对应一个inode，每个inode对应一个号码(int型), inode包含了文件的`元数据(metadata)`:
* 文件类型:  regular file, directory, pipe, socket等
* 所有者信息: user id，group id
* 权限信息: read, write, execute
* 时间戳: inode最后修改时间，文件内容最后修改时间，文件最后访问时间(`精确到纳秒`)
* Size: 文件字节数
* Blocks: 存储文件内容的blocks
* 链接数: 指向次inode的文件数量

注意: inode信息中不包括`文件名`信息

## 如何查看文件的inode

### ls -i命令
<img src="https://github.com/grearter/blog/blob/master/inode/ls.png" />

### stat命令
<img src="https://github.com/grearter/blog/blob/master/inode/stat.png" />

## inode存储在哪
当我们新建一个磁盘分区时，默认的该磁盘分为2个区域: 
* `inode区`: 存放inode信息
* `data区`: 存储文件内容数据

### 查看磁盘的inode使用情况
使用`df -i`命令:
<img src="https://github.com/grearter/blog/blob/master/inode/df.png" />
* Inodes列: inode区域大小
* IUsed列: inode区域使用空间
* IFree列: inodes区域剩余的空间

### 查看inode数据结构的大小
<img src="https://github.com/grearter/blog/blob/master/inode/inode_size.png" />
注意: 由于一个磁盘分区的`inode区域大小`在分区创建时已经确定, 若inode区域空间耗尽将无法创建新的文件(即使磁盘分区还有空间)

## inode链接数
`链接数(Links)`表示指向此inode的文件数量，linux允许多个文件指向同一个inode

### 什么时候`链接数`会变化
* 当删除文件(例如使用rm命令)时，inode的链接数`减1`，当链接数`减为0`时，此inode被回收
* 创建一个文件`硬链接(hard link)`时，inode的链接数`加1`

#### 硬链接(hard link)
使用ln命令创建一个文件的硬链接:
<img src="https://github.com/grearter/blog/blob/master/inode/hard_link.png" />

* 创建硬连接之后，a.txt 的链接数(Links)变成2
* a.txt 与 b.txt 指向了同一个inode(修改b.txt内容时a.txt内容也会变更)
* 仅删除a.txt 或 b.xt，此inode仍然存在(不会被系统回收)


### 软连接(symbolic link)
使用`ln -s <源文件> <链接文件>`来创建一个文件的软连接，软连接会创建新的inode，但文件的内容是源文件的路径(相当于windows的快捷方式)
<img src="https://github.com/grearter/blog/blob/master/inode/symbol_link.png" />
* 创建软链接之后，a.txt 的链接数不变
* 软链接文件的inode与源文件的inode不同
软链接文件与源文件为2个独立的文件，软链接文件仅仅是指向源文件的一个快捷方式而已！

## 目录的inode
在linux系统中，目录也是一种文件，保存了一系列`目录项`的列表，每个`目录项`包括:
* 文件名
* 文件名对应的inode号码


## FAQ
### Q1: 重命名一个文件inode是否会变更
A1: 对文件进行重命名(rename)操作，不会造成inode变化

### Q2: 移动(move)一个文件的位置，此文件inode是否会变更
A2: 若在同一分区内移动，不会改变文件的inode; 若从一个分区移动到另一分区，则inode信息会发生变化

### Q3: 进程打开一个文件的过程
1. 通过文件路径找到文件的inode号码
2. 通过inode号码在inode区域中找到inode信息
3. 根据inode信息找到文件所在的blocks，读取文件内容

### Q4: 若进程P1已经打开了文件A，此时对文件A进行重命名(rename)或移动(move)，进程P1是否会报错/抛出异常
A4: 不会。进程打开文件A之后，以inode来识别文件，而对文件rename或move不会对inode造成影响，因此进行P1仍可正常使用文件A

### Q5: 若进程P1已经打开了文件A，此时进程P2使用rm命令删除了文件A，进程P1是否会报错/抛出异常
A5: 不会。进程P2删除了文件A之后，文件A的链接数(Links)减为0，但由于此时进程P1仍然在使用文件A，所以操作系统不会回收文件A的inode，一直到进程P1关闭了文件A时，系统才会回收文件A。即当一个文件链接数为0且没有进程使用时，系统才会回收文件对应的inode与blocks。

### Q6: 如果通过inode来删除一个文件
A6: find YourTargetPath -inum InodeNum
