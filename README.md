# 自己写容器工具

使用的技术：
- linux namespace隔离技术
	- 挂载隔离  CLONE_NEWNS
	- nodename和domainname系统标记隔离   CLONE_NEWUTS
	- IPC进程通讯隔离  CLONE_NEWIPC
	- PID隔离   CLONE_NEWPID
	- 网络隔离  CLONE_NEWNET
	- User隔离  CLONE_NEWUSER
- linux cgroup 进程分组管理机制
- linux subsystem  资源控制模块，作用到cgroup上
	- blkio 块设备隔离
	- cpu cpu调度策略
	- cpuacct 统计cpu占用
	- cpuuset 设置进程能使用的cpu
	- devices 设备访问
	- freezer 对cgroup进程进行挂起恢复
	- memory 内存占用
	- net_cls cgroup网络包分类，用于区分出来对流量进行限流和监控
	- net_prio 设置 cgroup 中进程产生的网络流量的优先级
    - ns 这个 subsystem 比较特殊，它的作用是使 cgroup 中的进程在新的 Namespace fork新进程 CNEWNS ）时，创建出 个新的 cgroup ，这个 cgroup 包含新的 Namespace的进程
- hierarchy 用于组织cgroup为一个树结构
- 联合文件系统

