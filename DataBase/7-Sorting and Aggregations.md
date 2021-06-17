# Sorting and Aggregations

## 查询计划（query plan）

- 查询计划指的是指令，或者是数据库系统该如何执行一个给定查询的方式
- 查询计划被整理为一个树形结构或者是一个有向无环图，每个节点是一个operator
- 数据从树的叶子流向树根
- 根节点的输出是查询的结果

![image](https://user-images.githubusercontent.com/29897667/120910797-7a1bfe80-c6b4-11eb-8bcc-36f44f3577ea.png)

## 磁盘型数据库的operator问题

磁盘型数据库，各个operator产生的中间结果也没法完整放在内存中。

- 设计的operator算法要使用buffer pool
- 使用最大化顺序访问的算法

## 排序

- 在关系模型中，tuple间是没有指定的顺序的
- 但是很多query希望以有序的方式访问tuple
  - 根据属性去重（DISTINCT）
  - Aggregation（GROUP BY）
  - 加载大量数据，将其排序并自下至上地构建B+Tree索引，会更快
- 如果数据可以加载在内存中，则DBMS可以使用标准排序算法，如快速排序；若数据不可以全部加载在内存中，则DBMS需要使用外部排序算法，这样需要磁盘，且倾向于使用顺序访问而非随机访问。

## 外部归并排序（External Merge Sort）

### 排序过程

假如有B个buffer page可用，总的数据量为N个buffer page

- **Sorting**：在mem中对数据块（runs）进行排序，将排序后的runs写入到disk上的子文件。一个run的大小为B个buffer page。
- **Merging**：将排好序的子文件合并成更大的有序文件，一次merge B-1个runs，用一个page存储排序结果，排满一个page时，将其写到disk中。

![image](https://user-images.githubusercontent.com/29897667/120917051-af8a1180-c6df-11eb-8318-7ce76aa4e5c7.png)

### Double Buffering Optimization

在DBMS执行当前run时，Prefetch下一个run或是page到mem中，这样可以减少因为IO请求而带来的等待时间，且利用了顺序读更快的特性。

可以利用多线程，以便进行sort或merge，一边prefetch pages。

### 复用B+Tree

如果想要排序的key上建有B+Tree，则可以复用B+Tree索引，而不必使用External Merge Sort。

即从叶子节点最左侧向右遍历，这样无排序计算代价，且所有disk访问都是有序的。

## Aggregation

**Aggregation** 将一个或多个元组的值合并为单个标量值。实现聚合的方法有两种：**Sorting** 和 **Hashing** 。

### Sorting Aggregation

- DBMS首先在`GROUP BY` 修饰的key上将tuples排序。若数据可以全部加载如mem，则使用内部排序算法，否则使用external merge sort。
- 然后DBMS顺序扫描排好序的tuples去计算Aggregation，operator的输出在key上是有序的

### Hashing Aggregation

在DBMS扫描table时填充临时hash table。对每条record，检查在hash table中是否已存在该record（`DISTINCT`、`GROUP BY`）

#### External Hashing Aggregation

第一步：**Partition**

- 使用hash function *h1*， 基于hash key将tuples分partition
- 当partition满时将其写入disk（通过buffer pool manager）
- 这样保证具有相同key的record在一个partition当中，不需要去别的partition寻找是否具有相同key的record（同一个partition中保存的key可能不同）
- 假设有B个buffer，使用B-1个用于partition，剩下一个用于保存输入

![1623940916901](C:\Users\XutongLi\AppData\Roaming\Typora\typora-user-images\1623940916901.png)

第二步：**ReHash**：

使用hash table进行总结，将其压缩为要计算结果所需的最少信息。

- 对disk上的每个partition，读取其pages到memory中，并建立in-memory hash table（使用第二个hash function *h2*，*h1!=h2*）
- 然后遍历此hash table的每个bucket，将匹配的tuple集合起来，以计算聚合（假设每个partition都能放在mem中）

 ![image](https://user-images.githubusercontent.com/29897667/122419132-4b0f6200-cfbd-11eb-9408-54550ab4f071.png)

在 **ReHash** 期间，DBMS存储`(GroupByKey, RunningValue)`对来计算Aggregation。`RunningValue` 的值取决于Aggregation Function。向hash table中插入新的tuple时：

- 如果匹配到已有的`GroupByKey`，则更新`RunningValue`
- 否则插入`(GroupByKey, RunningValue)` 对

![1623941995854](C:\Users\XutongLi\AppData\Roaming\Typora\typora-user-images\1623941995854.png)

