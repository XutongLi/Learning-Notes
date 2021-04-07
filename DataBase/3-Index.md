# Index

## 1. Search Key

用于在DB文件中查找记录的属性或属性集称为 **搜索码 (Search key)** 。如果文件上有多个索引，那么它就有多个search key。

## 2. 基本索引类型

- **顺序索引**：基于值的顺序排序
- **散列索引**：基于将值平均分布到若干散列桶中。一个值所属的散列桶是由hash function决定的

## 3. 顺序索引

### 3.1. 聚集索引和非聚集索引

**聚集索引 (Clustering Index)**

- 文件中记录按照某个key的指定的顺序存储，则该key对应的索引为聚集索引。
- 聚集索引又称为主索引。
- 聚集索引通常建立在主码上，但建立在别的key上也可以。
- 在key上有聚集索引的文件称为索引顺序文件。

**非聚集索引 (Nonclustering Index)**

- key指定的顺序与文件中记录的物理顺序不同。
- 非聚集索引又称为辅助索引。
- 必须为稠密索引，每个key都对应一个索引项。

### 3.2 稠密索引和稀疏索引

**索引项 (Index Entry)** 是由一个key和指向对应的一条或多条record的指针构成的。指向record的指针为 `page_id + offset`。

#### 3.2.1. 稠密索引 (Dense Index)

- 文件中每个key都有一个索引项
- 在稠密聚集索引当中，索引项包括key和指向对应第一条记录record的指针。其余record顺序地存储在第一条record之后。
- 在稠密非聚集索引中，索引项必须包括指向所有对应记录的指针列表。

#### 3.2.2. 稀疏索引 (Sparse Index)

- 只为key的部分值建立索引项
- 只有索引是聚集索引时才能使用稀疏索引，因为需要顺序遍历查找
- 索引项包括key和指向对应第一条记录record的指针。其余record顺序地存储在第一条record之后
- 为了寻找某key对应的record，先定位到小于等于该key的最大key Index，再顺序遍历找到该key对应的record

#### 3.2.3. 对比

通过稠密索引可以更快地定位到一条record，但是稀疏索引空间和维护开销小。

**为每个page建立一个索引项的稀疏索引是一个比较好的折中**：因为处理数据库查询的开销主要由把page从disk读到memory中的时间决定。一旦将page放入memory中，其顺序扫描的时间就可忽略。使用这样的稀疏索引可以定位到包含所要查询记录的块，这样，只要记录不在溢出块里，就能使访问次数最小且索引尽可能小。

### 3.3 外层索引

当索引过大而不能放在内存当中时，需要存储在disk中，当需要时，从disk中复制page到memory中。如果索引占有b个page，在索引文件上使用二分法搜索定位索引项也需要访问 `logb` 个磁盘块。

可以在原始的内层索引上构造稀疏的外层索引。因为索引是有序的，这使得外层索引可以是稀疏的。当寻找一个record时，首先通过二分法在外层索引上找到对应索引项，在从这个索引项中指针指向的索引块中找到记录对应的索引项。

假设外层索引全部存储在memory中，那么当使用多级索引时，一次查询只需要读取一个page（内部索引）。索引比使用二分法在内部索引上搜索要少很多disk IO次数。

### 3.4 多个属性上的索引

key不是单个属性，而是一个属性列表。key值按照字典序排序。

如 `(dep_name, salary)` ，这一key是由系名和教师工资连接而成。

该复合key可用来进行的高效查询：

```sql
// 可高效查询，在两个属性上进行点查询
select ID	from instructor
where dept_name = "Finance" and salary = 8000
// 可高效查询，在第一个属性上指定等值条件，在第二个属性上指定范围查询（对应于搜索属性上的一个范围查询）
select ID	from instructor
where dept_name = "Finance" and salary < 8000
// 可高效查询，等价于(Finance, -∞)到(Finance, +∞)的范围查询
select ID	from instructor
where dept_name = "Finance"
// 不可高效查询，因记录存在于不同磁盘块，会导致大量IO（不对应于搜索属性上的一个范围查询）
select ID	from instructor
where dept_name < "Finance" and salary < 8000
```

[多级索引（复合索引）](<https://www.pianshen.com/article/9523749466/>)

[最左匹配原则](<https://www.cnblogs.com/lanqi/p/10282279.html>)

## 4. 位图索引

**位图 (bitmap)** 就是位的一个简单数组。

对于一个关系r，它的一个属性A只能取很少的一些值，则可在此属性上建立 **位图索引**。为该属性的每个取值建立一个位图，位图大小为记录数N。对每个记录编号，若记录i取值v，则将v值对应的位图上的第i位置1。该属性所有取值的位图共同构成位图索引。

**作用**：

- 优化在多个码上的选择操作：如 `select * from r where gender='f' and level='L2'` ，将gender属性和level属性的指定值上位图进行交运算，结果中值为1的位即为所要查询到的记录。
- 统计满足所给定条件的记录数：统计交操作后值为1的位数，可以在不访问表的情况下得到满足条件的记录数
- 存在位图：删除记录会造成存储空隙，移动记录来填充间隙代价又大。所以可以建立存在位图。该位图中如果第i位的值为0，表示记录i不存在，否则为1。
- 压缩B+Tree叶节点：B+Tree叶节点可能存储拥有某key的记录列表。记录以记录号 (page_id+offset) 的形式存储。此时可以用位图来标识记录，占用空间更小。

## 5. SQL中的索引操作

```sql
create index dept_index on instructor(dept_name)	// 定义索引
create unique index dept_idx on instructor(dept_name)	//在候选码上定义索引
drop index dept_idx 	// 删除索引
```















