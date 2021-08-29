# **Timestamp Ordering Concurrency Control**

## 1. TO Concurrency Control

TO并发协议是一种乐观并发协议，用于DBMS假定txn间的冲突很少的时候。它使用timestamp去确定txns的serializability order，而不是在读写对象前加lock。

![image](https://user-images.githubusercontent.com/29897667/126910288-a2b81516-9e00-408c-8585-81ae8e68ba7b.png)

## 2. Basic Timestamp Ordering Protocol (Basic TO)

Txns在没有lock的情况下读写obj。

![image](https://user-images.githubusercontent.com/29897667/126941791-f388eec4-feef-40d3-9f86-a3d3a168fc78.png)

![image](https://user-images.githubusercontent.com/29897667/126941861-198b20af-5da1-4230-92dc-d8aba1611435.png)

![image](https://user-images.githubusercontent.com/29897667/126941927-4fe91dd7-5592-4486-b8c4-67d3f24238c2.png)

**优化：Thomas Write Rule** 

![image](https://user-images.githubusercontent.com/29897667/126942509-7a4e6f2d-5f22-45a6-8fe5-f2ad10036386.png)

![image](https://user-images.githubusercontent.com/29897667/126942658-f6841665-72e0-4a6c-a053-46b35e14229e.png)

**recoverable**：一个事务只有当它所依赖数据的对应事务都已提交的情况下，再进行提交，这样的schedule是recoverable。

Basic TO Protocol下，DBMS不能保证txn读取到从故障恢复后将恢复的数据。

**缺点**：

- 将数据拷贝到txn的workspace以及更新timestamp有很大的开销
- 长的txn可能会因为和短的txn冲突而不断重启，出现“饥饿”现象
- 在高并发系统中会遭遇timestamp分配的瓶颈

## 3. Optimistic Concurrency Control (OCC)

### 3.1. 介绍

如果txns之间的冲突很少，且大多数txn都比较短，那么可以更乐观地假设txns之间没有conflicts。

OCC在冲突很少时表现很好，因为冲突少时加lock会带来很大的开销。

DBMS为每个txn创建一个**private workspace**：

- 读的obj都拷贝到workspace中
- 修改也应用到workspace

当txn提交时，DBMS比较workspace write sets去判断是否和别的txn冲突。如果没有冲突，write set被应用到db中。

### 3.2. OCC Transaction Phases

- **Read Phase**: 拷贝read sets到private workspace中，并在workspace中生成write sets
- **Validation Phase**: 当一个txn commit的时候，检查它是否和别的txn冲突
- **Write Phase**: 如果验证成功， 应用private change到DB。否则abort或restart txn（必须在一个受latch保护的critical section中执行Validation和Write phase）

### 3.3. Validation Phase

在此阶段，DBMS检查一个txn是否和别的txn冲突。DBMS需要保证仅有serializable schedule。

DBMS在txn进入验证阶段时给它安排timestamp。

Ti检查是否和别的txn有RW或WW冲突，并确认所有的冲突都是一个方向上的（从更老的txn到新的txn）。DBMS检查正在提交的txn和别的正在running的txn的timestamp顺序。

![image](https://user-images.githubusercontent.com/29897667/126999732-f3b42c37-36f0-4ef3-a9f4-947c5905724e.png)

### 3.4. 存在的问题

- 将数据拷贝到txn的private workspace开销比较大
- Validation和write phase由于要获取内存中的latch，因此有瓶颈
- OCC的abort回退的操作比别的协议更多，因为它是在commit时才abort
- timestamp分配也有瓶颈（需要在锁保护下递增）

## 4. Partition-based Timestamp Ordering

OCC中，因为同时只能有一个txn位于Validation Phase进行检查，因此需要latch，这样会造成性能低下。

一种优化方法是将DB分隔成不相交的 **partition**，只检查在相同partition上执行的txn的冲突。

每个partition通过一个lock保护，txn按照到达时间被安排timestamp。每个txn在启动前在partition的queue中排队。

**执行过程**：

- queue中第一个txn获取partition的lock
- txn在获取所有它将访问到的partition的lock后开始启动
- txn可以读写它lock的任何partion。如果一个txn尝试去访问没有lock的partition，那么它将abort+restart，接着获取这两个partition的lock，然后开始执行该txn

![image](https://user-images.githubusercontent.com/29897667/127036743-fa8abfd8-c2d1-46a6-aafd-a924f43fba9c.png)

**优点**：

- 所有的更新都是inplace的，不需要进行拷贝。会维护一个in-mem buffer，在txn abort的时候撤销修改
- 如果DBMS在txn启动前直到它要访问的partition，且大多数txns只访问单个partition，则这种方法很快

### 5. 幻读

在有插入或删除操作时，可能会有 **幻读** 的问题：

![image](https://user-images.githubusercontent.com/29897667/127037428-3c72071f-9247-49c8-bf4c-578faa5532ee.png)

原因是T1仅lock了原有的tuples。

**解决方法**：

1. 如果在status字段上有稠密索引，则txn可lock含有status='lit'的index page；如果不存在status='lit'的records，则加gap lock，即lock这个record将要在的page
2. 如果没有合适的index，txn必须要有table中所有page的lock或者table本身的lock

## 6. Isolation Levels

**Serializability** 可使程序员忽视并发问题，但是它带来更小的并发度并限制了性能。所以引入弱等级的一致性取提升可伸缩性。

**隔离级别** 控制txn暴露于其他并发txn操作的程度：

- **Serializable** 
  - 强制事务串行执行，这样多个事务互不干扰，不会出现并发一致性问题
  - 不发生幻读、所有读可重复、无脏读
  - 预先获取lock：strict 2PL + index locks
- **repeatable reads** 
  - 保证在同一个txn中多次读取同一数据的结果是一样的

  - 幻读可能发生

  - 同上，无index locks

  - （三级封锁协议：在二级的基础上，要求读取数据 A 时必须加 S 锁，直到事务结束了才能释放 S 锁。

    可以解决不可重复读的问题，因为读 A 时，其它事务不能对 A 加 X 锁，从而避免了在读的期间数据发生改变。即S2PL）
- **read committed** 
  - 一个txn的修改在提交之前对其他事务不可见

  - 幻读和不可重复读可能发生

  - 同上，S lock马上被释放

  - （二级封锁协议：在一级的基础上，要求读取数据 A 时必须加 S 锁，读取完马上释放 S 锁。

    可以解决读脏数据问题，因为如果一个事务在对数据 A 进行修改，根据 1 级封锁协议，会加 X 锁，那么就不能再加 S 锁了，也就是不会读入数据。）
- **read uncommitted** 
  - txn中的修改，即使没有提交，对其他txn也是可见的

  - 所有冲突都可能发生

  - 同上，没有S lock

  - （一级封锁协议：事务 T 要修改数据 A 时必须加 X 锁，直到 T 结束才释放锁。

    可以解决丢失修改问题，因为不能同时有两个事务对同一个数据进行修改，那么事务的修改就不会被覆盖。）