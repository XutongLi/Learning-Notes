# Join Algorithm

本节内容专注于两个table的 **inner equijoin** 算法，这些算法也支持别的类型的join。

在query plan中，会希望更小的表是left table（outer table， loop join中的概念）

## 1. Output

对于R表中的tuple r 和S表中的tuple s在join属性上匹配0时，join操作会组合r和s到一个新的output tuple中。

### Data

将outer table和inner table中属性的值直接复制到output tuple中。

这种方法的优点是查询计划中的未来操作永远不需要回到原表来获取更多数据。缺点是需要更多mem来存放整个元祖。

### Record IDs 

DBMS只复制匹配到的tuple的join key属性以及record id。这种方式适用于列存，因为DBMS不需要赋值query不需要的数据。称为**late materialization**

## 2. Cost分析

分析不同join算法的不同主要基于disk IO次数。包括从disk读取数据的IO以及将中间结果写入disk的IO。

> M pages in table R, m tuples total
>
> N pages in table S, n tuples total

## 3. Nested Loop Join

### 3.1. Simple Nested Loop Join

![image](https://user-images.githubusercontent.com/29897667/123422250-84188980-d5f0-11eb-8169-107e96162c73.png)

- **outer table** 是在外层for循环中的table，**inner table**是在内层for循环中的table。
- **join** 符号左侧侧表为outer table，右侧为inner table。
- ![image](https://user-images.githubusercontent.com/29897667/123422327-a14d5800-d5f0-11eb-88bc-f87c00befe91.png)
- 应该把更小的表当做outer table
- 尽可能将outer table缓存在mem中，如果可以的话，给inner table应用index

### 3.2. Block Nested Loop Join

![image](https://user-images.githubusercontent.com/29897667/123423463-2ab15a00-d5f2-11eb-9a25-659233f0643a.png)

- 对于outer table中的每个block，取inner table中的每个block并比较两个block中的tuples
- 让outer table和inner table中的每个page上的tuple进行join
- 这种算法降低了disk IO，因为是对每个outer table的block遍历inner table，而不是对每个outer table的tuple遍历inner table。
- ![image](https://user-images.githubusercontent.com/29897667/123463434-197e4280-d61e-11eb-9638-6e06f494fdf6.png)
- 如果有B个buffer可以使用
  - 则使用B-2个buffer进行outer table的扫描（对于outer table，用尽可能多的buffer保存它）
  - 1个buffer用于inner table
  - 1个buffer用于存储output 
  - ![image](https://user-images.githubusercontent.com/29897667/123463339-fbb0dd80-d61d-11eb-8d0a-85ecdeca0cd4.png)

### 3.3. Index Nested Loop Join

![image](https://user-images.githubusercontent.com/29897667/123555670-b1fcfa00-d7b9-11eb-8fb7-174bfcd2ab34.png)

- inner table使用索引，outer table不使用
- 可以是原有的索引，也可以是临时索引
- ![image](https://user-images.githubusercontent.com/29897667/123555689-dce74e00-d7b9-11eb-91f1-3d5dc7c38273.png)

### 3.4. 注意事项

- 更小的表作为outer table
- 尽可能将outer table缓存在mem中，以减少disk IO
- inner table使用index

## 4. Sort-Merge Join

![image](https://user-images.githubusercontent.com/29897667/123681566-90b11200-d87c-11eb-9bd7-98f21a744d1b.png)

**阶段一：Sort**

- 在join keys上排序两个table
- 若mem够用，则使用快排；否则使用external merge sort

**阶段二：Merge**

- 每个table一个cursor，遍历两个有序表，将join key一样的tuples作为output
- 可能需要进行回溯，只有inner table需要回溯（在outer table的join key重复，且和inner table的join key匹配时）

**Cost**：

![image](https://user-images.githubusercontent.com/29897667/123683113-69f3db00-d87e-11eb-90a7-4e90759a46a2.png)

**Worst case**：outer table和inner table在join key上都相等。但在现实数据库中这不会发生。

**适用场景**：这个算法适用于在一个或两个table都在join attributes上有序的情况。

## 5. Hash Join

**核心思想**

- 使用hash table把tuples基于join attributes分成更小的桶，这样会减少比较计算（不需要逐个tuple进行比较）。Hash Join只可以用于在complete join key上的equi-join。
- 如果R表中tuple r和S表中tuple s满足join条件，它们在join attribute上是相等的。使用一个hash function将r哈希到桶*ri*、将s哈希到桶*si*，之后只需要比较 *ri* 和 *si* 当中的tuples。

### 5.1. Basic Hash Join

![image](https://user-images.githubusercontent.com/29897667/123821512-0aed9f00-d92e-11eb-88da-088819857c5c.png)

**阶段一：Build**

- 在outer table上的join attributes上使用hash function *h1* 建立hash table（如果DBMS知道outer table的大小，则可以使用static hash table，否则使用dynamic hash table或者使用overflow page，如果hash table溢出到磁盘上，就需要进行大量随机IO）。key为join attribute，value形式包括
  - **Full Tuple**：避免在比较value时查询outer table的tuple；但是使用更多存储空间
  - **Tuple Identifier**：适用于列存，因为DBMS不需要从disk取不需要的数据

**阶段二：Probe**

- 遍历inner table，使用 *h1* 将每个tuple映射到对应bucket里，并通过值比较寻找匹配的outer table的tuple。

**Probe Phase Optimization**

- 在build phase，填充hash table时，建立Bloom Filter，它占用空间小，可以存储在内存中。
- 在probe phase，将bloom filter传到join里面。对hash table检测前，先检测bloom filter
  - 若检测为F，则表示outer table不含该key，无需查hash table，从而减少了磁盘IO
  - 若检测为T，去hash table中检测
- 有时称为 **边路信息传递**

### 5.2. Grace Hash Join / Hybrid Hash Join

为了处理不能再内存中进行的hash join（数据无法放在内存中）

![image](https://user-images.githubusercontent.com/29897667/124394663-057bc480-dd33-11eb-8a1e-7621d8ee8ad7.png)

**阶段一：Build**

- （普通Hash Join中，只为outer table构建hash table，然后对inner table检测是否有复合join条件的tuple，然后join）
- 对outer table和inner table都使用相同的hash function *h1* 创建hash table。hash table的buckets可以写到disk中
  - 如果一个bucket不能放在mem中，使用 **recursive partition**，即使用另一个hash function *h2* 将这个bucket分为子bucket

**阶段二：Probe**

- 检查两个hash table的对应bucket，在两个pages中使用nested loop将join attr对应的tuple进行匹配。这些page在mem中，所以减少了随机disk IO。

**Cost**：

![image](https://user-images.githubusercontent.com/29897667/124394712-35c36300-dd33-11eb-9e60-8d9de29fbe17.png)

## Summery

![image](https://user-images.githubusercontent.com/29897667/124395198-5b516c00-dd35-11eb-83f6-6f4785c45177.png)


## Extra: Bloom Filters

概率型数据结构（bitmap），用于回答近似成员查询（判断key是否存在于集合中），可插入和查找，但不能删除。

- False negatives 不会出现（即判断为F，实际为T）
- False positive 可能出现（即判断为T，但实际为F）

**步骤**：

- **Insert**：使用k个hash function，将hash结果对应的filter中bit置为1
- **Lookup**：检查k个hash function 哈希后的bit位是不是都是1：若都是1，则判断为T，否则为F