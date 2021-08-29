# Index Concurrency Control

## 1. Concurrency Control

需要使多线程安全地访问索引结构，以利用额外的CPU核并且隐藏disk IO停顿。多数DB支持多线程，也有的DB不支持多线程， 如Redis，即使只有单个线程访问它的数据结构，依然可以获得好的性能。

**Concurrency Control Protocal** 是DBMS用来保证对于共享对象进行并发操作的正确性的方法。

## 2. Locks VS Latches

**Locks**：

- 多个事务执行时保护DB中的逻辑内容，逻辑内容可以是tuple、tuple的集合、或一张表、或是数据库等等
- 在事务执行期间持有Lock
- 可以支持将对象回滚至操作前的状态。如A给B转账，B未到账时系统崩溃，需要在系统恢复撤销未完成的操作。

**Latches**：

- 在多线程间保护DBMS In-Memory Data Structure的临界区（critical section）
- 在操作执行期间持有Latch
- 不需要支持回滚

![image](https://user-images.githubusercontent.com/29897667/109762557-be1de400-7c2b-11eb-9449-f859a4e36b93.png)

## 3. Latches

### 3.1. Modes

- **Read Mode**：多个线程可同时读一个共享对象。一个线程可以在其他线程持有read latch时获得read latch。
- **Write Mode**：仅一个线程可以修改对象。一个线程可以在其他线程不持有任何latch时获得write latch。

### 3.2. Implementation

可以用于实现latch的底层原语 是通过现代CPU提供的 atomic `compare-and-swap (CAS)` 指令实现的。一个线程可以检查内存（进程地址空间中）位置中的值，判断其是否为某指定值，若是，则CPU将该值修改；否则不修改。

有不同的方法可以实现latch，这些方法在实现复杂性和性能之间存在trade-off。这些 `test-and-set` 操作具有原子性，在一个线程check了value但未update时，其他线程不能update这个value。

#### 3.2.1. Blocking OS Mutex

使用OS内置的 `mutex` 作为latch。如 std::mutex。代价昂贵且无可扩展性（每次lock/unlock调用大概25ns）

![image](https://user-images.githubusercontent.com/29897667/109777218-6d63b680-7c3e-11eb-8408-fa8f16b2fc57.png)

Linux内核使用 `Futex(fast user-space mutex)` ，它有两个组成部分：在user-space中的一个spin latch、OS级的mutex。DBMS优先获取user-space latch，若获取失败，进程切换入内核态并尝试获取更昂贵的mutex。如果获取mutex也失败，则线程通知OS它因申请锁而阻塞，OS将其放在一个等待队列中，等待调度器调度。

[Futex](https://www.jianshu.com/p/d17a6152740c)

#### 3.2.2. Test-and-Set Spin Latch (TAS)

使用CPU提供的CAS指令对进程空间中的值进行判断，不会像mutex一样进行线程状态切换（用户态->内核态），速度很快。如果CAS失败，则通过一个while循环尝试继续更新。循环中也可实现为满足某条件时挂起让其他线程执行或退出。

但是无可扩展性，对cache不友好，循环会消耗CPU资源。如std::atomic\<T>

![image](https://user-images.githubusercontent.com/29897667/109778378-b5370d80-7c3f-11eb-9c32-62ca0c7f794f.png)

#### 3.2.3. Reader-Writer Latch

如果应用程序读操作很多，那么使用读写锁可以很好地提高性能，因为它允许并发的读操作访问共享对象。读写锁在spin latch的基础上实现。

读写锁记录多少个线程正在持有锁，通过管理不同的队列来跟踪不同类型的latch有哪些线程正在等待获取。

**starvation**：一个线程无法访问到一个共享对象，如由于加读锁的线程太多，导致申请写锁的线程无法获取写锁。

DBMS需要管理读队列和写队列去避免starvation。

## 4. Hash Table Latching

因为所有thread都沿着相同的方向遍历，且每次只访问一个slot/page，所以对于hash table的并发控制比较容易，且不会死锁。

如果要改变表的大小，与要给整张表加个锁。

根据latch粒度大小的不同，Hash Table Latch的类型可分为：

- **Page(Bucket) Latch**：每个page都有自己的读写锁。它减少了并行性，因为一个page同时只能被一个thread访问，但如果是线程内顺序访问page内slot会比较快。

  ![image](https://user-images.githubusercontent.com/29897667/109833174-0107a800-7c7c-11eb-8c91-26a8ab8c0853.png)

- **Slot Latch**：page中每个slot都有自己的锁。它提升了并行性，因为一个page可以同时被多个thread访问。但是提高了存储和计算开销。DBMS可以使用简单模式的latch（如spin latch）去减少元数据和计算开销。

  ![1614787004557](C:\Users\XutongLi\AppData\Roaming\Typora\typora-user-images\1614787004557.png)

## 5. B+Tree Latching

### 5.1. 基本方法

要使得多个线程可以同时读写一个B+Tree，需要解决两种问题：

- 多个线程同时修改一个节点的内容
- 一个线程正在遍历树，而另一个线程正在拆分或合并节点

**Safe Node** 是更新时确定不会被分裂或合并的节点，它的状态满足：

- 非全满（保证插入时不会被分裂）
- 多于半满（保证删除时不会被合并）

**Latch Coupling** 方法用于保证多线程安全地访问一棵B+Tree，其步骤为：

- **Find**：从根节点开始向下遍历，重复地：获取子节点的读锁，释放父节点的读锁
- **Insert/Delete**：从根节点开始向下遍历，获取途径所有节点的写锁。一旦一个子节点被锁，检查它是否safe，若safe，则释放它所有祖先节点的锁。
- 在向下遍历B+Tree的过程中，线程会用一个stack保存一路上持有的所有latch

### 5.2. 优化

每次insert/Delete操作都需要在根节点加写锁，写锁是exclusive的，降低了并行性，这会造成很大的性能瓶颈。

假设大多数操作不会造成叶节点的分裂和合并（真实DB中一个node大约有16KB，会存有很多key），于是insert/delete操作先获取节点读锁，到叶节点时获取写锁，**如果叶节点不是safe的**，释放当前持有的所有锁，从根节点开始重新获取写锁。

（即先尝试使用乐观锁）

### 5.3 Leaf Node Scans

不考虑兄弟指针的情况下，B+Tree只会被从上至下地访问，这种情况下不会发生死锁。

考虑兄弟指针时，在叶节点遍历上就有了两个方向（从左至右和从右至左）。Index Latch不支持deadlock detection 和 deadlock avoidance。

于是解决这一问题的唯一方法是通过编码规则。获取兄弟节点的锁时必须支持 `no-wait` 模式，即thread试图获取叶节点上的锁，但该锁不可获取时，那么它将立即中止其操作（释放所有持有的锁），并从根节点重新启动该操作。

### 5.4 Delayed Parent Updates

处理overflow时的一种额外优化手段

发生overflow时要拆分节点，需要修改三个节点：

- 更新被拆的节点-
- 添加新节点
- 父节点或祖父节点中添加新节点中的key用作路标

**Delayed Parent Updates** 是处理overflow时的一种额外优化手段，当一个叶节点溢出时，延迟更新它的父节点。（这样就不用重头开始并拿着写锁一路往下遍历了，只需要更新这棵树中全局信息表中的一点内容）

![1614843015064](C:\Users\XutongLi\AppData\Roaming\Typora\typora-user-images\1614843015064.png)

如上图所示，插入25时，将31分裂到了新节点中，此时不会重新获取写锁更新C，而是C的修改记录在一个全局信息表中。当后续线程持有C的写锁时，再完成这一修改。

## 6. 本节课核心思想

- 沿一个方向加锁时不会出现死锁
- 当出现死锁时中断