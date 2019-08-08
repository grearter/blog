# inode

## 什么是inode

在Linux系统中，每个文件(包括目录)都对应一个inode，inode包含了文件的`元数据(metadata)`:
* 文件类型： regular file, directory, pipe, socket等
* 所有者信息：user id，group id
* 权限信息：read, write, execute
* 时间戳：inode最后修改时间，文件内容最后修改时间，文件最后访问时间(`精确到纳秒`)
* Size: 文件字节数
* Blocks: 存储文件内容的blocks
* 链接数: 指向次inode的文件数量

## 如何查看文件的inode

### ls -i命令
<img src="https://github.com/grearter/blog/blob/master/inode/ls.png" />

### stat命令
<img src="https://github.com/grearter/blog/blob/master/inode/stat.png" />
