# Lec6 - Fault Tolenrance Raft 2

![](./img/1.png)

为什么S2的index12中的term 4的entry可被丢弃？

少数派服务器所在分区中，无论是哪个term的leader对外发送的AppendEntries，都不会得到来自多数派分区服务器们的确认（网络割裂了）。因此，前一个leader并不能决定这个Clien端的命令是否被提交了；也不能去执行它，更不能将它提交到应用状态上；也永远不会发送一个肯定的回复给Client端（指日志落地）；因为多数派服务器没有提交Client这个命令，Client得不到响应, 就没有理由去相信这条命令已经被（leader）执行了。

**persistent**：每当添加log entry 或 改变 current term或 vote 时都要进行持久化写入。（可以在响应或发送RPC之前再写入、一次持久化100条entry、在发送AE之前要持久化log entry）

**snapshot**：state machine 状态要比对应的 log 小，所以当log大小超过某个阈值后，将一个log point对应的执行状态进行备份。

解决建立快照时，follower无法赶上leader的问题：install snapshot RPC



***

