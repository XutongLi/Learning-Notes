# Logging Protocols + Schemes

## 1. Crash Recovery

恢复算法是确保DB一致性、事务原子性以及持久性的技术。

任何故障恢复算法包括 **两部分**：

-  txn执行时记录可以保证DB故障后恢复的信息
- 故障后使DB恢复到一个保证atomicity、consistency以及durability的状态的操作

故障恢复算法的两个关键原语（并不是所有故障恢复算法都有）：

- **UNDO**：撤销未提交事务在disk上的修改
- **REDO**：重新执行已提交事务为被持久化的修改

## 2. Buffer Pool Management Policies

### 2.1. Steal Policy

DBMS是否允许未提交事务修改非易失性存储中被已提交事务修改过的对象（事务是否可以往disk上写未提交的修改）

### 2.2. Force Policy

在事务提交前，DBMS确保事务所做的修改反映在非易失性存储性。

### 2.3. 不同策略效果

- **STEAL/NO-FORCE**：性能最好，也带来更复杂的日志设计和恢复处理。需要UNDO和REDO。
- **STEAL/FORCE**：仅需要UNDO能力。
- **NO-STEAL/NO-FORCE**：仅需要REDO能力。
- **NO-STEAL/FORCE**：什么也不需要，性能最差，需要大量内存。

## 3. Shadow Paging

DBMS为page维护master和shadow：

- **master**：仅包含已提交事务做的修改
- **shadow**：保存未提交事务做的修改

事务在shadow上写。当事务提交时，原子性地切换shadow为master。

**策略**： **NO-STEAL/FORCE**

**实现**：

- 将DB pages组织为树状，root是一个disk page
- root指向master，在shadow上进行更新
- 事务提交时，将root指向shadow，并将修改后的root page落盘，然后修改内存中指针，即交换master和shadow

**UNDO**：移除shadow page

**REDO**：不需要

**缺点**：

- 提交时的刷盘是随机IO，性能低下
- 会有磁盘碎片，浪费空间

## 4. Write-Ahead Logging

维护一个记录对DB修改的log

- log保存在持久化存储中
- log包含进行UNDO和REDO的足够信息

**策略**：**STEAL/NO-FORCE**

**WAL协议**：

- DBMS将事务的log records存在buffer pool的page中
- DBMS需要在被修改的objects落盘之前将相关的log落盘
- 对于一个事务，在它的log records落盘之后，它才可以提交
- 事务开始时log记录`<BEGIN>` 
- 事务结束时，log中记录 `<COMMIT>`，确保事务提交前所有的log records落盘
- 每个log entry包含object的一些信息：
  - Transaction ID
  - Object ID
  - Before Value (UNDO)
  - After Value (REDO)
- 可以使用 **group commit** 去分批flush多个事务的log，以提高性能

**优点**：写入是顺序IO，性能好

**缺点**：恢复时要重新执行事务，耗时长

## 5. Logging Schemes

### 5.1. Physical Logging

- 记录磁盘物理位置上所做的字节级变化
- 如：git diff

### 5.2. Logical Logging

- 记录 txn high-level 的操作
- 如事务中的 UPDATE、DELETE、INSERT
- 相对于physical logging，写入的数据更少
- 缺点：
  - 恢复时要重新执行每个txn，因此使用时间更多
  - 如果有并行txn，很难实现logical logging，因为很难确定一个query已经修改了DB的哪一部分

### 5.3. Physiological Logging

- log records针对单个page，但不需要指定page的数据组织方式
- 使用最多

![image](https://user-images.githubusercontent.com/29897667/128632942-1518e963-ff4f-4f73-b8a5-f88b92d13f89.png)

## 6. CheckPoints

log会逐渐变多，占用很大空间，且恢复时需要执行整个log，导致需要很长恢复时间。

因此DBMS可以定时生成 **checkpoint**，它flush 所有的page。

**多久生成一个checkpoint**：太频繁会导致大量IO从而降低性能、太少会导致恢复时需要太多时间以及需要更多空间存储。

**Blocking CheckPoint Implementation**：

- DBMS停止接收新的事务，并等待所有活跃事务完成
- 刷新所有log和dirty data page到disk
- 在log中写入`<CHECKPOINT>`，并将该log刷入disk

![image](https://user-images.githubusercontent.com/29897667/128633197-c57d937b-44ff-4f73-8b00-1bf3290889e5.png)



