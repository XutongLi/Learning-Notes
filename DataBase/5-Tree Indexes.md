

# Tree Indexes

## 1. Table Indexes

**table indexes** 是表中属性的子集的一个 **副本** ，它们被以一种高效的方式进行存储，以允许我们进行高效的查找（相比循序遍历）。

由于index是副本，所以DBMS会确保表内容和index是同步的。如果改变了表中一个tuple，也会将这个修改反应在index上。

对于DB的index数量的选择存在trade-off，需考虑存储成本（buffer pool中和磁盘中）和维护成本（每次插入和更新都要修改index）。

> 若使用Hash Index，则必须使用完整的key，无法进行部分查找或范围查找，如>、like。
>
> 在表的一个属性上可以建立多种类型的Index，query plan会按照查询的具体类型来确定使用的Index。

## 2. B+Tree Basic

### 2.1 B+Tree 介绍

[二叉搜索树、平衡二叉树、B-Tree、B+Tree 区别](<http://www.liuzk.com/410.html>)

**B+Tree** 是 self-balancing的树形结构，它会保证数据的有序性，支持 **O(logn)** ($log_{\lceil n/2 \rceil}n$) 复杂度的查找、顺序访问、插入和删除。

相比B-Tree，B+Tree有一个优势是：当遍历到B+Tree底部时，可以沿着叶子节点进行循序扫描。

### 2.2 B+Tree 特性

- B+Tree是一种多路查找树（M-way search tree），这意味着在树中的每个节点处，它可以通过M条不同的路线到达其他节点。
- 它是 **perfect balanced** 的，每个叶节点都在相同的高度；对树进行修改时，该数据结构会始终保持平衡性
- 每个中间节点有 `ceil(M/2) - M` 个子女
- 每个节点都至少是半满的情况：`ceil((M-1)/2) <= #keys in one node <= M-1` （M为节点的度）
- 每个有 `k` 个key的中间节点都有 `k+1` 个非空子节点

### 2.3. B+Tree Node

#### 2.3.1. Overview

![1614455576737](https://user-images.githubusercontent.com/29897667/113731912-bb7f4480-972b-11eb-8c8a-c3f20faf8abe.png)

- B+Tree每个节点都是由key-value对构成的数组，node中的key以某种规则排序
- B+Tree 叶节点含有Sibling Pointers（兄弟指针）用于循序遍历，中间节点没有
- 中间节点是\<node指针|key>结合体，key是在任意属性上构建的索引
- 开始用一个给定的key进行搜索时，通过key的值确定其搜索路线
- 叶子节点是\<value|key>结合体，value可能是record id，也可能是tuple

#### 2.3.2. 叶节点具体结构

![image](https://user-images.githubusercontent.com/29897667/109398813-15952900-797a-11eb-94e6-e4d4fdc948b2.png)

- level表示叶节点在第几层
- slots表示空闲slot数量
- prev、next指针分别指向前一个和后一个叶节点
- key和value分开存储是因为value大小可能不同，且key连续放在一起方便都放入CPU cache中进行二分查找
- value中可能存储：
  - Record ID：即指向key对应tuple的指针
  - Tuple Data：将tuple直接存入叶节点，二级索引必须将record id存储为它们的value。

### 2.4. B-Tree VS B+Tree

- B-Tree每个节点都存有key和value，它更节省空间，因为不会存储重复key
- B+Tree仅在叶节点中存储value，中间节点只作为路标，所以会有重复的key
- B+Tree中，当将一个key删除时，可能将其保存在中间节点（取决于是否需要重组），如果要查找其他key，还可以通过这条路线往下找
- 不使用B-Tree原因：使用多个线程进行更新操作时代价大，对一个内部节点的修改需要被向上和向下传播，这个在并发操作下是需要被保护的；而B+Tree只对叶节点修改，只可能有向上传播

### 2.5. B+Tree Insert

- 通过中间节点的路标（数组中的key值）找到要插入的叶节点
- 准备将data entry有序地插入到叶节点中
  - 如果叶节点 `L` 未满，直接插入
  - 否则将 `L` 从中间key分裂为 `L` 和 `L2` ，均匀地重新分配entry，向上复制中间key，并将指向 `L2` 的索引项插入到父节点中（需递归操作，因为父节点也可能已满）

[B+Tree Insert Demo](https://www.cs.usfca.edu/~galles/visualization/BPlusTree.html)

### 2.6. B+Tree Delete

- 从根节点开始查找，找到entry所在的叶节点 `L` ，删除该entry
- 若 `L` 中的key数量满足半满，即 `ceil(M/2)-1 <= #keys <= M-1` ，则完成
- 否则：
  - 尝试重新分布，从有相同父节点的兄弟节点借一个key给 `L` 
  - 若重新分布失败，合并 `L` 和有相同父节点的兄弟节点
- 对于中间节点，至少要有 $\lceil n/2 \rceil$ 个指针；对于叶节点，至少要有 $\lceil (n-1)/2 \rceil$ 个key。

## 3. 聚簇索引

**Clustered Indexes**：将数据与索引存储在了一起，聚簇索引具有唯一性。DBMS会保证对page中tuple的物理布局记性匹配排序。

有的数据库一定会使用聚簇索引，如果table不包含primary key，那么DBMS会自动生成一个隐藏的row id primary key。

好处：如果使用聚簇索引，那我们知道tuple存储的顺序和主键是一致的，只需要遍历某个叶节点下包含的一部分pages，就可以找到所要的数据。否则tuples可能存储在不同page上，需要进行大量随机IO。

## 4. B+Tree Design Choices

### 4.1. Node Size

存储设备越慢，B+Tree节点应该越大。

因为如果跳到不同节点的随机IO速度更快，就可以使节点更小；否则应使节点内顺序访问尽可能多。

- HDD：1MB
- SSD：10KB
- In-memory：512B

### 4.2. Merge Threshold

理论上当一个节点中的key不足半满时就可能进行merge操作。但有的DBMS会延迟merge操作，这样可以减少重新组织的数量（如先删除1再插入1，可以省去merge操作）。

允许节点不半满的情况存在并定时重新组织整个树会更好。

### 4.3. Variable Length Keys

处理变长key的方法：

- node中存储指向作为key的tuple的属性的指针。因为tuple存储在disk中，要做大量磁盘IO，所以不使用此方法。

- 使用变长Node。因为buffer pool中page大小固定，所以维护变长Node需要进行复杂内存管理，此方法也不使用。

- 使用0或NULL将key填充至定长。有的DBMS会使用，耗费空间大。

- Key Map / Indirection。在Node中维护一个offset数组，存储key+value在Node中的offset。使用最多。

  ![image](https://user-images.githubusercontent.com/29897667/109520916-ad645580-7ae7-11eb-8686-a7c6f695af75.png)

### 4.4. Non-Unique Indexes

**对于一个节点中key重复的处理**：

- 在节点中存储重复的key

  ![image](https://user-images.githubusercontent.com/29897667/109540881-0ccd6000-7afe-11eb-827e-e3671596ce30.png)

  查询时找到最左的key，在向右遍历，直到找到所有包含该key的records。

- 只存储key一次，但是将value存储在一个list中

  ![image](https://user-images.githubusercontent.com/29897667/109540947-22db2080-7afe-11eb-8a9d-8d8e7e48b2f9.png)

**对于整个B+Tree中重复key的处理**：

- **Append Record Id**：将key与tuple的record id (page_id + offset) 组合成为新的唯一的key。查找时，相当于使用新key的前缀来查找。实际中主要使用这种方法。

  ![image](https://user-images.githubusercontent.com/29897667/109553472-0b0b9880-7b0e-11eb-8dd0-cede7771f140.png)
  ![image](https://user-images.githubusercontent.com/29897667/109553525-18c11e00-7b0e-11eb-8775-fe2e5d798d43.png)

- **Overflow Lead Nodes**：允许为叶节点扩展溢出节点。这种方法实现复杂。

  ![image](https://user-images.githubusercontent.com/29897667/109553592-2d9db180-7b0e-11eb-93f7-a42ee9c9280c.png)

### 4.5. Intra-Node Search

几种在节点中进行搜索的方式：

- 顺序扫描
- 二分查找，需要key在node中是有序的
- 差值法：根据已知的key的分布规律估计key的大约位置，以此位置作为扫描的起始位置

### 4.6. 辅助索引和记录重定位

一些文件组织（如B+Tree文件组织）可能会改变记录的位置，即使该记录并未被更新（如B+Tree节点的分裂和合并）。若非聚集索引叶节点中存储了指向记录的指针（记录的位置），则记录位置改变时，会引发非聚集索引更新而导致的disk IO。

解决这一问题的方法是：在非聚集索引的叶节点中，存储聚集索引的key。避免了record位置更新带来的非聚集索引节点的IO操作。但是查找非聚集索引后还需要查找一次聚集索引。

## 5. B+Tree Optimazations

### 5.1. Prefix Compression

**前缀压缩**：在一个叶节点当中，若key是有序的，则它们很有可能有相同前缀。所以可以将该前缀提取出来，只存储这些key在该前缀后的部分。可以节省大量空间。

（在字符串属性上建索引可以使用该优化，因为字符串不定长，且字符串比较大的话会减少节点的出度）

![image](https://user-images.githubusercontent.com/29897667/109545778-3be6d000-7b04-11eb-93ba-41f3fd501e67.png)

### 5.2. Suffix Truncation

**后缀截断**：因为中间节点中的key只用于寻路，所以不需要存储完整的key。可在中间节点中只存储保证正确寻路的前缀。

![image](https://user-images.githubusercontent.com/29897667/109547014-dd225600-7b05-11eb-806a-3402bd7d3a10.png)

![image](https://user-images.githubusercontent.com/29897667/109547052-e6132780-7b05-11eb-82c9-3f5e4906b432.png)

### 5.3. Bulk Insert

**大量插入**：有时需要将大量数据导入数据库，逐个插入key同时构建index效率很低下，因为要进行大量的合并。可以将所有数据插入后，再将key排序，自底向上地构建索引。

### 5.4. Pointer Swizzling

一个节点中保存的是其子节点或兄弟节点的page id而非指针，所以需要先通过buffer pool来获取page在内存中的地址（指针）。

如果page被pin在buffer pool中，那么可以在节点中保存指针，这样避免了page table的查找（page table查找需要加锁）。

## 6. 其他索引使用

### 6.1. 隐式索引 (Implicit Indexes)

在创建主键或生命唯一性约束时，DBMS会隐式创建Index去实现完整性约束。但是创建外键（外键是另一张表的唯一性索引）时不会隐式创建Index。

![image](https://user-images.githubusercontent.com/29897667/109603917-9a8c6800-7b5d-11eb-896d-a2a4a6e9b19e.png)

### 6.2. 部分索引 (Partial Indexes)

可以在表的部分tuple上构建索引，这么做可以减少索引大小，且可以减少维护索引的开销。

部分索引最常用的场景是对于日期范围划分索引，如对每个月建立一个索引。

使用部分索引可以避免一堆不需要的数据取污染buffer pool。

![image](https://user-images.githubusercontent.com/29897667/109639973-d5f15b80-7b8a-11eb-958e-a009f31bcdb3.png)

多数数据库支持。

### 6.3 覆盖索引 (Covering Indexes)

处理查询所需的所有字段都能在索引本身中找到，这样的索引称为覆盖索引。DBMS不需要去查询tuple。（不需要去回表）

此举减少了buffer pool中的锁争用。

![image](https://user-images.githubusercontent.com/29897667/109644651-b52c0480-7b90-11eb-8493-2cf48b374fda.png)

少数数据库支持。

### 6.4 Index Include Columns

在索引中嵌入额外的属性，以支持仅访问索引的查询。

这些额外属性仅存储在B+Tree叶节点中，且不作为key的一部分。

![image](https://user-images.githubusercontent.com/29897667/109653110-92ebb400-7b9b-11eb-8f89-336f7112626b.png)

### 6.5 函数式/表达式索引 (Functional/Expression Indexes)

索引不一定需要以key在表中的方式来存储key。

可以使用表达式来建立索引。

![image](https://user-images.githubusercontent.com/29897667/109654869-adbf2800-7b9d-11eb-9efc-e86900c12b72.png)

## 7. Trie

### 7.1 简介

使用key的 **digit**（bit、byte或其他）表示来逐个检查前缀，而不是比较整个键。又称为 **Digital Search Tree**、**Prefix Tree**。

从根节点到某个节点，其路径上digit组成的key，即为该节点对应的key。

两个具有相同前缀的key，它们在Trie上游相同的起始路径。

![image](https://user-images.githubusercontent.com/29897667/109666527-3a6fe300-7baa-11eb-8c42-a3a07ee8f049.png)

### 7.2 特征

- 形状取决于key的组成和长度，与key的插入顺序无关，不要求进行重新平衡的操作
- 所有操作复杂度为 `O(k)` ，`k` 为key的长度
- 不存储完整key，由从根节点到叶节点的路径来表示key

- 点查询Trie更快，但对于扫描来说，Trie比B+Tree慢很多，因为要进行很多回溯操作

## 8. Radix Tree

也叫作 **Patrucia Tree**。对Trie做垂直压缩，忽略所有仅有一个child相传的路径。

有可能出现false positive（如对于前缀ha，hat和hair都匹配），所以tree搜索出Record_id后需要通过该id找到tuple，与key作对比以验证key是否真匹配。

下图一为Trie，图二为对齐垂直压缩后的Radix Tree：

![image](https://user-images.githubusercontent.com/29897667/109700078-4f5d6e00-7bcc-11eb-98ee-078156316c25.png)

![image](https://user-images.githubusercontent.com/29897667/109700139-600de400-7bcc-11eb-87fc-d6c1ea7de29f.png)

## 9. Inverted Index

**倒排索引** 被用来映射一个word和目标属性中含有该word的record（即试图找的是某个属性中的一个子元素）。

一般用于keyword search。（Hash Table适用于point search，B+Tree适用于range search）

**支持的查询类型**：

- 词组查询：查找含有一组给定顺序的words的records
- 近似查询：查找两个words在对方n个words内的records
- 通配符查询：寻找符合某种pattern的words所在的records

**存储什么**：Inverted Index至少需要存储每条记录中包含的words（用标点符号分隔）。它还可以包括额外的信息，如词频、位置和其他元数据。

**何时更新**：更新Inverted Index是开销非常大的，所以大多数DBMS都会维护辅助数据结构来分段更新，然后批量更新索引。


