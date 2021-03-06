# 进程状态

## 概述
进程是当前正在执行的计算机程序的实例。在进程的生命周期中, 一个进程会经历多种状态。

## 进程类型
在linux系统之中, 进程有以下几种类型:
* 用户进程(user process)
* 守护进程(daemon process)
* 内核进程(kernel process)

### 用户进程
系统中的大多数进程都是用户用户进程, 用户进程一般是指由普通用户创建的、在用户态执行的进程。

### 守护进程
守护进程是一种在后台(background)运行的应用程序, 通常用于管理某种正在进行的服务。<br/>
守护进程可能会监听&处理传入的请求, 例如, httpd守护程序舰艇&处理http请求。 <br/>
守护进程可能会随着时间的推移自行启动活动, 例如, crond守护程序自动在预设时间启动cron任务。
出于安全考虑, 非必需场景, 守护进程不应使用root权限来启动, 而是使用特定的用户来启动, 例如syslogd进程使用syslog用户启动, ntpd进程使用ntp用户启动。
守护进程随系统启动, 并一致持续到系统关闭。

### 内核进程
内核进程仅在内核空间中执行, 与守护进程非常类似, 主要区别在于内核进程可以完全访问内核数据结构(比用户空间中运行的守护程序进程权限更高)。
但内核进程也不如守护进程灵活, 更改内核进程可能需要重新编译内核。

## 进程状态
* 新建(create), 一般通过`system()`或`fork() + exec()`系统调用来创建一个新的进程
* 就绪 (ready),进程具备运行条件(资源得到满足或等待的事件己经发生), 等待系统分配处理器以便运行
* 运行(running), 进程占有CPU正在执行指令
* 等待(waiting), 进程正在等待事件或资源
	* `可中断等待`: 可通过信号中断
	* `不可中断等待`: 在任何情况下都不能被中断(如设备驱动等待磁盘或网络IO), 不可中断等待状态的进程不会立即处理(handle)信号, 而是进程退出等待状态时, 才会(延迟)处理在等待状态收到的信号
* 停止(stoped), 进程处于停止状态, 一般来说进程收到SIGSTOP信号进入停止状态
* 退出(terminated)

## 进程状态切换
* 新建->就绪: 当一个进程被创建后, 如果当前系统中有足够的内存, 则进程进入`就绪`状态；
* 就绪->运行: 进程调度器(按照调度策略为进程分配CPU资源)在某一时刻为进程分配CPU资源, 此时进程进入执行状态；
* 运行->就绪: 当CPU时间片耗尽, 调度器又会将剥夺进程的CPU资源, 此时进程重新进入`就绪`状态；
* 运行->等待: 当进程需要等待某项资源(如IO操作)时, 只有在获得等待的资源后才能继续执行, 则进程进入`等待`状态；
* 等待->就绪: 当等待的资源得到满足后, 处于`阻塞`状态的进程转换到`就绪`状态;
* 运行->退出: 进程执行完成自动退出或收到终止信号(如SIGTERM);
<img src="https://github.com/grearter/blog/blob/master/linux/process_state/process_state.png" />

## 查看进程状态
使用`ps`命令或`top`命令, 可以查看每个进程当前的状态信息。
* R: 运行状态(running) 或 就绪状态(runnable)
* S: 可中断的睡眠状态(interruptible sleep)
* D: 不可终端的睡眠状态(uninterruptible sleep)
* T: 停止(stopped), 进程为结束但处于停止状态(可以理解为特殊的等待状态)
* Z: 僵死(zombie)
<img src="https://github.com/grearter/blog/blob/master/linux/process_state/ps.png" />

### 例1: 状态R(就绪或运行状态)
```c
#include <stdio.h>

// run_state.c

int main() {
	printf("process start\n");
	while (1) {
		;
	}
}
```
编译并运行: `gcc run_state.c -o run_state && ./run_state`, 使用ps命令查看进程状态:
<img src="https://github.com/grearter/blog/blob/master/linux/process_state/r_state.png" />

### 例2: 状态S(可中断的等待状态, 等待IO)
```c
#include <stdio.h>

// interruptible sleep

int main() {
	printf("process start\n");

	int a = 0;

	scanf("%d", &a);

	return 0;
}
```
编译并运行: `gcc interruptible_sleep.c -o interruptible_sleep && ./interruptible_sleep`, 使用ps命令查看进程状态:
<img src="https://github.com/grearter/blog/blob/master/linux/process_state/interrunptible_sleep.png" />

### 例3: 状态S(可中断的等待状态, 等待事件)
```c
#include <stdio.h>
#include <unistd.h>

// interruptible_sleep2.c

int main() {
	printf("process start\n");
	sleep(60); // sleep 60s
	return 0;
}
```
编译并运行: `gcc interruptible_sleep2.c -o interruptible_sleep2 && ./interruptible_sleep2`, 使用ps命令查看进程状态:
<img src="https://github.com/grearter/blog/blob/master/linux/process_state/interrunptible_sleep2.png" />

### 例4: 状态T(停止状态)
```c
#include <stdio.h>

// stop_state.c

int main() {
	printf("process start\n");
	// some code to ensure process do not exit
	while (1) {
		;
	}

	return 0;
}
```
编译并运行: `gcc stop_state.c -o stop_state && ./stop_state`, 使用ps命令查看进程状态:
<img src="https://github.com/grearter/blog/blob/master/linux/process_state/stop_state.png" />


### 例5: 状态Z(僵死状态)
```c
#include <stdio.h>
#include <sys/types.h>
#include <unistd.h>

// zombie_state.c

int main() {
	printf("process start\n");

	pid_t ret_pid = fork();

	if (ret_pid < 0) {
		// fork failed
		return -1;
	} else if (ret_pid == 0) {
		// child process
		printf("i am child process, pid=%d\n", getpid());
	} else {
		// parent process
		printf("i am parent process, pid=%d\n", getpid());
		while (1) {
			;
		}	
	}
}
```
编译并运行: `gcc zombie_state.c -o zombie_state && ./zombie_state`, 使用ps命令查看进程状态:
<img src="https://github.com/grearter/blog/blob/master/linux/process_state/zombie_state.png" />
