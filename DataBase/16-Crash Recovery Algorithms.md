# Crash Recovery Algorithms

## 1. ARIES

Algorithm for Recovery and Isolation Exploiting Semantics

**ARIES 恢复算法的三个关键概念**：

- **Write Ahead Logging**：在DB将dirty page写入disk前，该page上相关修改的log会被写入disk（STEAL+NO-FORCE）
- **Repeating History During Redo**：重启时，重新执行操作到crash前的特定状态
- **Logging Changes During Undo**：记录撤销操作到log上，以在重复故障时确保操作不会重复

## 2. WAL Records

每个log record有一个全局唯一的 **log sequence number (LSN)**。添加日志的时候，会通过一个单调递增的counter来给每条日志记录分配LSN。

![image](https://user-images.githubusercontent.com/29897667/128642801-b4f3123b-4aba-4a3f-9ba8-f243f6fde548.png)

![image](https://user-images.githubusercontent.com/29897667/128642807-70c6c5d5-c085-42c9-a92f-eb70dc72cd31.png)

- 当事务修改page时更新 **pageLSN**
- 当DBMS将WAL buffer写入disk时更新 **flushedLSN**

## 3. Normal Execution

每个事务调用一系列读和写操作，最后时commit或abort。

### 3.1. Transaction Commit

当一个事务要commit的时候：

- 写 `<COMMIt>` record 到 mem的log buffer中
- DBMS flush 到commit entry的所有的log records到disk（顺序的、同步的）
- 一旦 `<COMMIT>` record存于disk上，DBMS返回ACK给应用，告诉应用事务已提交
- 此时事务内部实际未完全完成，依然要维护一些内部元数据描述事务实际做了哪些事情（如修改bookkeeping表，这个表是用来标识所有活跃的事务的）。
- 之后DBMS会写 `TXN-END`record到 log buffer中，标示该事务已经完全完成，不会再有和其相关的records
- `TXN-END` record不需要被立即刷盘

### 3.2. Transaction Abort

Abort 事务是一种特殊的UNDO，它只针对单一事务。

为log records额外添加一个字段 `prevLSN`，指向事务前一个LSN。DBMS通过这种方式建立事物中操作的反向链表，方便UNDO时找到事务的操作

引入新类型的log record **compensation log record (CLR)** 。CLR 记录了事物中需要UNDO的操作的反向操作。它的字段时update record的各字段加 `undoNext pointer`。CLR不需要被UNDO。abort时，DBMS会直接通知app，而不会等到CLR落盘。CLR的作用是使得事务正常执行过程可以执行UNDO（区别于故障恢复时的UNDO）。

![image](https://user-images.githubusercontent.com/29897667/128762361-2e327d37-9aec-4ec8-aee2-be2099a1a20e.png)

abort一个事务时：

- DBMS首先添加 `<ABORT>` 到 log buffer中
- 逆序地undo事务的更新操作，对于每个undo的update record，添加其对应的CLR
- 通过执行CLR来恢复旧值
- 在事务所有操作被撤销后，写 `TXN-END`到 log buffer

## 4. Checkpointing

DBMS周期性地创建CheckPoints，此时会将所有dirty page刷盘。这么做可以减小log大小，减少故障恢复所用时间。

### 4.1. Blocking Checkpoints

DBMS暂停所有的事务确保此CheckPoint可以生成一个一致性的快照：

- 停止所有新事务的开始
- 等待活跃事务执行完
- 将所有dirty page刷盘
- 容易实现，但是性能很低

### 4.2. Slightly Better Blocking CheckPoints

DBMS不需要等待活跃事务执行结束，但是DBMS会记录checkpoint开始时的系统状态

- 停止新事务的开始
- 在DBMS生成checkpoint的时候暂停事务执行

#### 4.2.1. Active Transaction Table (ATT)

- ATT中保存正在活跃的事务，当事务commit或abort后事务会被从ATT中移除。
- 对每个事务，ATT中保存：
  - *transactionId*：唯一的事务id
  - *status*：Running、Commiting、Undo Candidate
  - *lastLSN*：事务的最近的LSN

#### 4.2.2. Dirty Page Table (DPT)

DPT中记录了buffer pool中被未提交事务修改的dirty page的相关信息。一个dirty page对应一个entry，保存了 *recLSN*，是第一个使得page dirty的record的LSN。

![image](https://user-images.githubusercontent.com/29897667/128984092-b8f9e0b0-dbac-4f9a-8a76-79a87c92bb36.png)

### 4.3. Fuzzy CheckPoints

DBMS在CheckPoints时将dirty page落盘时，允许活跃事务继续执行。这是ARIES使用的协议。

添加新的log records来追踪CheckPoints的边界：

- `CHECKPOINT-BEGIN`：指示CheckPoint的开始（在checkpoint成功完成后，CHECKPOINT-BEGIN的LSN会写入disk的**MasterRecord**上，即MasterRecord保存了最后一次CheckPoint的begin log）
- `CHECKPOINT-END`：包含ATT + DPT（在CheckPoint开始后启动的事务不会加在ATT里）

## 5. ARIES Recovery

分为三个阶段：

![image](https://user-images.githubusercontent.com/29897667/129018650-71986d06-ce40-4f88-bf1e-0b8ea4713e7a.png)

![image](https://user-images.githubusercontent.com/29897667/129018751-2a53fec0-84a1-4899-acbe-3726c46a1978.png)

### 5.1. Analysis Phase

![image](https://user-images.githubusercontent.com/29897667/129018927-7b937288-778f-4e9f-90d5-8183d8039282.png)

![image](https://user-images.githubusercontent.com/29897667/129018961-5f3c26e5-07b0-477a-b88c-9037c1b5b188.png)

### 5.2. Redo Phase

![image](https://user-images.githubusercontent.com/29897667/129019552-ba9d37f9-df39-4349-8e94-e928541e9511.png)

### 5.3. Undo Phase

![image](https://user-images.githubusercontent.com/29897667/129019689-4f73c30c-5545-4b6e-9fb1-f0ff58f6ce6a.png)

### 5.4. Other issues

![image](https://user-images.githubusercontent.com/29897667/129022135-f138eb28-3041-4680-a556-7685a46883a8.png)



![image](https://user-images.githubusercontent.com/29897667/129022175-4ad7a330-c946-4a42-9bbc-93068d2c4ae8.png)









