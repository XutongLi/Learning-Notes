# Distributed OLTP Database Systems

  ## 1. Two Phase Commit

[MySQL 两阶段提交](https://blog.csdn.net/JesseYoung/article/details/37970271)

[对分布式事务及两阶段提交、三阶段提交的理解](https://blog.csdn.net/weixin_34365417/article/details/85816639?utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromMachineLearnPai2%7Edefault-1.control&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromMachineLearnPai2%7Edefault-1.control)

各节点会将每个阶段发生的事情以及他们所接收并发送的消息记录到WAL中。

**coordinator挂掉**

- participant会判断如何处理这种问题。

- 最简单的做法是：如果coordinator挂掉了，假设事务中止了，只需回滚事务所做的修改。但可以让participant意识到coordinator挂掉了，但是自己事务还在执行，于是可以变成新的coordinator。然后通过投票情况确定是否提交事务。

**participant挂掉**

- coordinator超时时间内收不到participant的ACK，会认为收到了abort

## 2. CAP 理论

![image](https://user-images.githubusercontent.com/29897667/129472140-a27ccbc9-d132-4847-a4c5-c2d81184d319.png)

**Consistency**：线性一致性，即强一致性。在所有节点上看到的数据都是一致的

**Availability**：任意给定时间，可以访问任何节点并获取系统中的任何数据（但不保证获取的数据为最新数据）

**Partition Tolerant**：如果丢失消息（网络故障、及其故障），仍然可以处理并得到想要的任何响应

NoSQL：AP（不停止更新，提供节点恢复连接后解决冲突的机制）

newSQL或传统的OLTP：CP或CA（停止更新，直到多数节点恢复连接）

