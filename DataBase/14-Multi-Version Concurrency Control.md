# Multi-Version Concurrency Control

## 1. MVCC

MVCC不仅是一个并发控制协议，更是影响了DBMS管理事务和DB的方式。

DBMS为一个**逻辑对象**维护多个**物理版本**：

- 当一个txn写obj的时候，DBMS创建这个obj的新的version（virtual snapshot）
- 当一个txn读一个obj的时候，DBMS读此txn开始时该obj的最新version（按照txn进入系统的时间分配timestamp）

**关键特征**：

- writer不会阻塞reader
- reader不会组设writer
- 只读事务可以在不加lock的情况下读一个一致性的snapshot
- timestamp用于确定可视性
- DBMS可通过MVCC支持**时序查询 (time-travel queries)**

![image](https://user-images.githubusercontent.com/29897667/127214650-3d88e297-f12b-46f7-9d8e-69d3ef278f3c.png)

## 2. Concurrency Control Protocol

- Timestamp Ordering
- Optimistic Concurrency
- Two-Phase Locking

## 3. Version Storage

需要弄清楚某个tuple的哪个版本是可见的

这是DBMS如何存储一个logical obj 的 physical version

![image](https://user-images.githubusercontent.com/29897667/127383417-d68d4ca1-1721-4dfa-8086-717ab86604c2.png)

### 3.1. Append-Only Storage

当创建数据的新版本时，只需复制该tuple的老版本，并将该副本作为表空间中的一个新物理tuple，并对其进行更新以及更新指针。

![image](https://user-images.githubusercontent.com/29897667/127378857-ac9351f9-a610-483a-a7a6-75ca2fb0fc89.png)

指针可以从旧到新，也可以从新到旧，取决于实现：

![image](https://user-images.githubusercontent.com/29897667/127378920-4cbdd550-d8e6-49a4-acc2-93ed5b2c5aeb.png)

### 3.2. Time-Travel Storage

master table保存tuple的最新版本，将老版本数据复制到Time-Travel table上

维护master version表中指向Time-Travel表的指针

![image](https://user-images.githubusercontent.com/29897667/127379656-373079b9-5a2d-43e1-a1c5-56be095e0212.png)

### 3.3. Delta Storage

维护对前一个版本所做的修改

![image](https://user-images.githubusercontent.com/29897667/127381987-f3b8fa42-9069-4895-b838-b7123c508e0c.png)

**优点**：写操作块；**缺点**：读旧版本数据慢

## 4. Garbage Collection

DBMS需要取移除可回收的、不再使用的物理版本。

### 4.1. Tuple Level Garbage Collection

- **Background Vacuuming**：独立线程定时地扫描表、检查过时的版本并清除它们（O2N、N2O下都可使用）
- **Cooperative Cleaning**：线程执行查询遍历version chain时，遇到旧版本的数据就将其回收（只可用于O2N）

### 4.2. Transaction Level 

每个txn维护自己的读写集，当一个txn结束时，garbage collector通过这个读写集来确定哪些tuple可以回收。

由DBMS来决定一个已完成的txn创建的版本何时不可见。

## 5. Index Management

**primary key index** 会指向version chain的head。txn update pkey被当作DELETE+INSERT。

![image](https://user-images.githubusercontent.com/29897667/127734860-44e5dff0-89fe-4db2-90ab-37af1e5e8f5f.png)

**secondary index** 的管理更复杂：

- **Approach 1 :Logical Pointers**

  - sec index中存放建立clu index的字段

  ![image](https://user-images.githubusercontent.com/29897667/127734927-512c1730-57bf-4a3b-8a93-407dc944512a.png)

  - 建立tuple id与physical address之间的映射，sec index中存tuple id。数据更新时修改映射表即可。

  ![image](https://user-images.githubusercontent.com/29897667/127734968-bfbb2db4-b4cb-418f-8034-9713f21f4323.png)

  

- **Approach 2 : Physical Pointers**

  - 使用version chain head的物理地址（这样更新数据时，所有的二级索引都要修改）

![image](https://user-images.githubusercontent.com/29897667/127734870-16772f21-6100-4cca-a9fc-6c2ff5055aa7.png)