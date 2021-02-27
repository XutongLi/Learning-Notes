# Database Storage

## 1. 存储层次
![image](https://user-images.githubusercontent.com/29897667/108463907-32ff2e80-72ba-11eb-8a6c-1617398d00ba.png)

- **Volatile Devices** 支持快速随机访问，以字节寻址。这意味着程序可跳转至任何字节地址并获取数据。以 *memory* 代称。
- **Non-Volatile** 顺序访问更快，它是块（block/page）寻址的，意味着如果要读取某个地址的数据，它会加载4KB的块到内存中，然后再取出数据。以 *Disk* 代称。
- 磁盘访问时间 = 寻道时间 + 旋转等待时间
  - 寻道时间：磁盘臂定位到正确磁道所需要的时间
  - 旋转等待时间：读写头到达所需的磁道后，等待访问的扇区出现在读写头下所花费的时间
- 提升磁盘块访问速度的方法：
  - Buffer Pool
  - Pre-fetching
  - 调度：按照使磁盘臂移动最短距离的顺序发出块访问请求
  - 文件组织：按照与预期的数据访问方式最接近的方式来组织磁盘上的块
- 现在有了一种新型存储设备 **non-volatile memory**，即 **持久化内存** ，它有着和DRAM差不多的速度并且可以持久化存储。

## 2. RAID

独立磁盘冗余阵列（Redundant Array of Independent Disk）

- 通过对读写操作并行执行（数据拆分）提高性能
  - 比特级拆分：将每个字节的第i位写到第i张磁盘上，这样磁盘吞吐率可以提高8倍
  - 块级拆分：将块拆分到多张磁盘。当读大文件时，块级拆分可以同时从n张磁盘上并性地读取n个块，从而获得对于大的读操作的高数据传输率；当读取单块数据时，其数据传输率与单个磁盘上的数据传输率相同，但是剩下的n-1张磁盘可以自由地进行操作。
- 通过在多张磁盘中存储冗余信息来提高可靠性

### 2.1 0级

块级拆分、没有冗余

![image](https://user-images.githubusercontent.com/29897667/109339745-0d1def00-78a3-11eb-8ba9-974fc31b5943.png)

### 2.2 1级

块级拆分、镜像

![image](https://user-images.githubusercontent.com/29897667/109339818-24f57300-78a3-11eb-9cf4-c05e30c829ab.png)

### 2.3 3级

比特级拆分，奇偶校验

![image](https://user-images.githubusercontent.com/29897667/109339837-30489e80-78a3-11eb-8433-9f69daef51cf.png)

如果任意一块硬盘发生错误，可以通过其余磁盘信息重建：对应磁盘对应扇区的对应位使用一个奇偶校验位，记录于多余的一个磁盘上。

3级相对于1级的好处：

- RAID3级对于多张常规磁盘只需要一个奇偶校验磁盘，而RAID1级对每张磁盘都需要一张磁盘镜像，因此RAID3级减少了磁盘开销
- 使用比特级拆分，对一个字节的读写散步在多张磁盘中，所以它使用N道数据拆分读写一个块的传输率是RAID1级的N倍。

### 2.4 5级

块级拆分、分布式奇偶校验

![1614364008197](C:\Users\XutongLi\AppData\Roaming\Typora\typora-user-images\1614364008197.png)

将数据和奇偶校验位分布到所有的N+1张磁盘当中，而不是在N张磁盘上存储数据并在一张磁盘上存储奇偶校验位。这样做，所有磁盘都能参与对读请求的服务，而不像RAID4级中存储奇偶校验位的磁盘不参与读操作，所以RAID5级增加了在一段给定时间中能处理的请求总数。

5级与1级相比（各有取舍）：

- 1级提供最好的写操作的性能，在数据库系统日志文件存储中应用广泛
- 5级相对于1级有较低的存储负载，但是写操作需要更高的时间开销；对于经常进行读操作而很少进行写操作的应用，5级是首选

5级与3级相比（3级不如5级）：

- 对于大量数据传输，5级与3级有同样好的数据传输率
- 对于小量数据传输，5级使用更少的磁盘。由于对于小量数据传输，磁盘访问时间占主要地位，所以并行读取并没有带来好处

### 2.5 6级

块级拆分、P+Q冗余

![1614364016775](C:\Users\XutongLi\AppData\Roaming\Typora\typora-user-images\1614364016775.png)

和RAID5相比，存储了额外的冗余信息，以应对多张磁盘发生故障的情况。

## 3. Why not use mmap
使用 **mmap** 可以将文件内容映射到进程地址空间，这样使得OS负责将block在disk和memory之间转移。但是在 **page fault** （进程在虚拟地址空间中找不到指定page时，会抛出page fault） 时会导致进程阻塞。

令DBMS管理block的转移有以下好处：
- 可以以正确的顺序将dirty page刷入disk中
- 可以进行特殊的page prefetching
- 可以使用特殊的Buffer替换策略
- 方便线程、进程调度

## 4. Disk-Oriented DBMS Overview
![image](https://user-images.githubusercontent.com/29897667/108473817-7e6d0900-72c9-11eb-9625-afe7ed759a14.png)
- 数据库存储在 **disk** 上，数据库文件中的数据组织为 **Page** 的形式，第一个page是 **directory page**。
- **page** 是固定大小的数据块，可以存储tuple、meta-data、indexes、log records 等
- 每个page都有一个单独的page_id
- 每个page都有一个header，存有page size、checksum、DBMS version、压缩信息等数据
- 为了对数据进行操作，DBMS需要将数据从disk拷贝到memory中，它使用 **buffer pool** 来管理数据在disk和memory之间的移动。
- DBMS有 **execution engine** 来执行查询（query），execution engine 向 buffer pool 请求一个指定的page，buffer pool会给execution engine一个指向该page在memory中位置的指针。

## 5. Page Storage Architecture
### 5.1 堆文件组织（Heap File Organization）
需要含有追踪哪些page存在且哪些page含有空间的元数据。
#### 5.1.1 Linked List
![image](https://user-images.githubusercontent.com/29897667/108604320-2a0c7b00-73e8-11eb-814c-ea0bc4b3dc3e.png)
- 维护一个header page，含有两个指针：free page list和data page list
- 每个page都追踪自己的free slots的数量
- 缺点：查找一个page时，需要进行顺序搜索
#### 5.1.2 Page Directory
![image](https://user-images.githubusercontent.com/29897667/108604489-44932400-73e9-11eb-9341-410701db40b8.png)
- DBMS维护一个特殊的page，用于追踪各data page的位置以及其中的free slots数量
- DBMS需要保证directory page和data page是同步的
## 6. Page Layout
数据在page中如何组织
### 6.1 Slotted-pages
![image](https://user-images.githubusercontent.com/29897667/108604741-c33c9100-73ea-11eb-85e9-bd69eb641bb0.png)
- page header中存储了已使用的slot数量、slot起始offset（slot倒序使用）以及一个slot array（存储了每个tuple的offset）
- 添加tuple时，slot array从前往后增加，data从后往前存储，当slot array和tuple data相遇时表示page被填满
- 删除tuple时，它所占用的空间被释放，page中在被删除tuple之前的tuple向后移动，header中的offset也被修改。移动记录的代价并不高，因为page的大小有限制。
- slotted-page方法要求没有指针指向tuple实际位置，而是指向slot array的相关位置。
- 当前DBMS中最常用的
### 6.2 Log-structured
![image](https://user-images.githubusercontent.com/29897667/108604954-e7e53880-73eb-11eb-9878-181101d956cc.png)
![image](https://user-images.githubusercontent.com/29897667/108604957-eddb1980-73eb-11eb-8295-fb5c30eccb76.png)
- DBMS只存储log record（db如何被修改，insert、update、delete）而非tuple
- 读取一个记录时，从后往前读log，并重建tuple
- 支持快速写，暂时慢速的读，适用于append-only的storage
- 为避免太多内容需要读，可建立索引帮助跳转到指定的log entry
- 定时进行压缩（如对于一个tuple有insert和update操作，可压缩为一个insert操作）
- Hbase、LevelDB、RocksDB等数据库使用

## 7. Tuple Layout
一个tuple即一串字节，DBMS负责将这这个字节串解释为属性类型和值。
### 7.1 定长记录
![image](https://user-images.githubusercontent.com/29897667/108605400-d7828d00-73ee-11eb-81c4-d83d2081a68e.png)
- 每个tuple都有一个header，记录visibility info（concurrency control）和对于Null Valuede Bit Map
- 属性以定义它们时的顺序存储
- 每个tuple有一个id（page_id+offset）
### 7.2 变长记录 
![image](https://user-images.githubusercontent.com/29897667/108605807-7f995580-73f1-11eb-8710-af6fba245358.png)
- 一个具有变长属性的记录表示通常具有两个部分：初始部分和变长属性的值。
- 对于定长属性，分配存储它们的值所需的字节数。
- 对于变长属性，如varchar类型，在记录的初始部分表示为一个对（偏移量，长度）值。偏移量表示该记录中该属性的开始位置，长度表示该属性的字节长度。
- 在记录的初始定长部分后，这些属性的值是连续存储的。
- NULL Bit Map用来记录哪个属性是空值。
- 对于超出page大小的value，可以使用overflow page，令tuple该属性指向这个overflow page。
- 对于大对象（图片、音频）的存储，将其存储在特殊文件中，将指向该对象的指针存在tuple中。
### 7.3 denormalize (prejoin)
![image](https://user-images.githubusercontent.com/29897667/108605725-0863c180-73f1-11eb-9b6b-f7adbe411b63.png)
![image](https://user-images.githubusercontent.com/29897667/108605735-131e5680-73f1-11eb-937d-28ffb13c2deb.png)
如果两个表是相关联的，DBMS会prejoin它们，所以这两个表会被存储在一个page中。这使得查询时DBMS只需要加载一个page就可以得到两者join后的结果。但是它会使得更新低效一点，因为对于每个tuple需要更多的空间。
### 7.4 Data Representation
- 存储在tuple中的数据类型主要有：integers、variable precision numbers、fixed point precision numbers、variable length values and dates/times。
- **Variable Precision Numbers**，如FLOAT、REAL等，运算速度比固定精度数字更快，因为CPU可以在它们上直接执行指令。但是会有进位错误。
- **Fixed Point Precision Numbers**，如NUMERIC、DECIMAL等，有灵活的精度和规模，通常存储在带有元数据的变长二进制表示中。使用它们不会有进位错误

## 8. System Catalogs
DBMS在自己的internal catalog中存储一些元数据，包括：
- Tables、columns、indexes、views
- Users、permissions
- Internal statistics

可通过查询 `INFORMATION_SCHEMA` catalog来获取数据库的信息。

## 9. Workload
![image](https://user-images.githubusercontent.com/29897667/108999505-5de3eb00-76dd-11eb-9fbc-4ea8aa462447.png)

### 9.1 OLTP
**Online Transaction Processing**

- 主要工作负载为针对少量记录的增删改，写多读少，并发访问多
- 具有实时性
- 主要用于基本的、日常的事务处理，如Amazon的购物操作、银行系统等
- 主要衡量标准为事务吞吐量，优化方法主要有：
  - 访问数据一般有热点，所以可以扩大内存容量，让buffer pool缓存更多的数据
  - 并发量高，CPU每秒要处理的请求也多，所以需要CPU处理能力更强
  - 与客户端交互数据量不大，但是频繁，所以要提高网络传输能力
### 9.2 OLAP
**Online Analytical Processing**

- 长时间运行的、更复杂的查询
- 读取表中大量记录，并发访问不多
- 一般用于分析从OLTP DB收集到的数据
- 对时间要求不严格
- 如Amazon会用OLTP存储用户购买行为，将所有数据存于OLAP对其进行数据分析
- 主要衡量标准为查询响应速度（QPS），优化方法有：
  - 单次访问数据量大，但是数据分布集中，所以需要有尽可能大的IO吞吐率，所以要选用吞吐率大的磁盘
  - OLAP系统每次运算过程较长，可以并行化（将一个任务，如select全表，分配到多个节点上），所以一般OLAP系统都是由多台主机构成一个集群，集群中主机间数据交互量大，所以需要提升集群内网络

### 9.3 HTAP

**Hybrid Transaction and Analytical Process**

- 核心为：如何在OLTP单一数据系统上，提供OLAP操作
- 数据不需要从操作型数据库导入到决策类系统
- 操作事务实时地对分析业务可见

## 10. Data Storage Models
### 10.1 N-ary Strorage Model (NSM)
- row store
- DBMS将一个tuple的属性顺序存储
- 适用于OLTP workload：插入和更新操作多
- 适用于对于整个tuple的查询
- 缺点：对于查询多数tuple中的少数几个属性的情况性能差，因为要读入大量无用属性到buffer pool
### 10.2 Decomposition Storage Model (DSM)
- column store
- DBMS连续存储所有tuple的同一个属性
- 适用于OLAP workload：查询表中多数tuple的某个属性。因为不需要读入无用属性
- 方便更好的压缩
- 缺点：对于整个tuple的查询、插入、更新、删除操作慢

