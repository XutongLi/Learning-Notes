# Concurrency Control Theory

## 1. Transactions

concurrency control和recovery基于TX ACID特性的概念。

**事务（transaction）**是通过在数据库中执行一系列操作来执行某种更高级的功能。它是DBMS中基本的变化单元。

操作交错执行导致的：

- 暂时性地不一致（不可避免）
- 永久性地不一致（不可以）

## 2. 定义

![image](https://user-images.githubusercontent.com/29897667/126307766-bf9e9caf-ab73-41f5-aa0c-d7f3044d9407.png)
![image](https://user-images.githubusercontent.com/29897667/126307820-34d63a8d-4fce-4892-8713-97ee0c5a57c5.png)

## 3. 事务的正确标准：ACID

### 3.1. 原子性（Atomicity）

事务中的操作要么都执行，要么都不执行。

#### 3.1.1 Logging

- DBMS记录所有的修改在log中，在tx abort的时候通过该log undo这些修改
- 在mem和disk中都维护undo records
- 通过logging可以将随机写入变成循序写入
- log还可以用来跟踪审计应用程序所做的每一件事

将文件记录在disk上，每次对数据库做修改，会把要

#### 3.1.2. Shadow Paging

- DBMS在事务执行前拷贝pages，事务在这些副本上执行。事务commit后再讲指针指向这个副本，并表示这个副本时数据的主版本
- 很少有数据库使用这种方法

### 3.2. 一致性（Consistency）

如果数据库是一致性的且事务是一致性的，当事务执行的时候，数据库的状态也是一致的。

一致性在这里指逻辑上的正确性。

数据库实际上是对现实世界中的某些概念进行建模。假设数据库的逻辑是正确的，即无需在意物理存储是什么样的，只要数据完整性和引用完整性等是正确的就行了，在此前提下当向数据库询问任何问题的时候，都会生成正确的结果。

对于单节点系统，一致性没有太多意义，对于分布式系统很重要。

[如何理解数据库事务中的一致性的概念](https://www.zhihu.com/question/31346392/answer/569142076)

#### 3.2.1. Database Consistency

- 数据库准确地表示它所建模的真实世界实体，并遵循完整性约束（Integrity Constraint）
  - 如一张学生表，约束年龄不能小于0
- 未来的事务会看到db中之前commit的事务的修改结果

#### 3.2.2. Transaction Consistency

- 如果数据库在执行事务前是一致的，那么事务执行后也是一致的
- 确保transaction consistency是应用的责任

### 3.3. 隔离性（Isolation）

同一时间有多个事务在执行时，事务间是彼此隔离开的。

#### 3.3.1. Concurrency Control

并发控制协议是DBMS如何决定正确交错执行多个Tx中的操作。我们希望db的最终状态与在单线程中按照顺序执行这些事务所得到的结果相同。

**两种协议类型**：

- **悲观协议**：不让并发冲突出现，即事务执行前加lock（两阶段锁协议）
- **乐观协议**：假设冲突很少发生，事务执行前不加lock，冲突后再解决（时间戳顺序协议）

##### 3.3.1.1. Schedule

DBMS执行操作的顺序称为 **execution schedule**。并发控制协议的目标是生成等价于一些序列执行的execution schedule：

- **Serial Schedule**：不使不同tx的操作交错的schedule
- **Equivalent Schedule**：对于任何数据库状态，执行第一个schedule的效果等价于执行第二个schedule
- **Serializable Schedule**：等价于tx的某些串行执行的schedule

**如何判断一个schedule是不是正确的**：如果schedule等价于某个顺序执行，则是correct的。

##### 3.3.1.2. Conflict

两个操作是**conflict**的，如果满足：

- 它们是不同Tx的操作
- 它们访问同一个object，且其中至少一个操作为write

conflict类型：

- **Read-Write Conflict (R-W)** — **Unrepeatable Reads (不可重复读)**：一个tx在读一个对象多次时不能得到相同的value

  ![image](https://user-images.githubusercontent.com/29897667/126539066-3477a892-563a-40cf-a946-b6ddd21d8670.png)

- **Write-Read Conflict (W-R)** — **Dirty Reads (脏读、读未提交数据)**：一个Tx看到了另一个Tx在commit前对一个对象的修改

  ![image](https://user-images.githubusercontent.com/29897667/126540699-192c9355-3954-49f7-a59a-b310dcc075e0.png)

- **Write-Write Conflict (W-W)** — **Lost Updates、Overwriting Uncommitted Data (覆盖写掉未提交的数据)**：一个Tx覆盖掉了另一个Tx未commit的修改

  ![image](https://user-images.githubusercontent.com/29897667/126541224-2adcb3d9-98ad-4fdf-9d5b-364a79708021.png)

##### 3.3.1.3. 序列化 (Conflict Serializability)

![image](https://user-images.githubusercontent.com/29897667/126549834-5e11f8e5-8200-4910-9262-fda50d103764.png)
![image](https://user-images.githubusercontent.com/29897667/126549891-24972556-9c00-4064-a2df-8890c91c5b12.png)
![image](https://user-images.githubusercontent.com/29897667/126549920-230115f6-ce3d-4add-915e-fe1ad6671529.png)

![image](https://user-images.githubusercontent.com/29897667/126549962-5be06b3a-ed5e-4699-b983-295b5c5e6fee.png)

### 3.4. 持久性（Durability）

如果所有的tx commit，所有的修改都会被永久保存。

DBMS使用logging或shadow paging确保所有的改动都是durable。

## 4. 总结

- 并发控制是自动的，DBMS自动地插入lock / unlock的请求和不同tx的schedule action
- 确保执行结果是等价于以某种顺序一个接一个地执行txs
- 事务可以帮助解释程序执行的正确性，即是否以正确的顺序执行写操作

