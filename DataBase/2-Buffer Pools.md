# Buffer Pools

## 1. Buffer Pool Organization

Buffer Pool 是内存中对于从disk读取的page的缓存，它是一个由固定大小page组成的数组。数组的每一个单元是 **frame**，当DBMS请求一个page时，这个page会被拷贝到一个frame中。

![image](https://user-images.githubusercontent.com/29897667/109200357-776b5c80-77db-11eb-93a7-7e9d7448b111.png)

Buffer Pool中维护的 **meta data** 有：

- **Page Table**：内存中的hash table，映射page_id和buffer pool中的frame位置，用于追踪哪些page正在buffer pool中
- **Dirty-flag**：当线程修改一个page时会set dirty flag，向storage manager指明了这个page必须被写回disk
- **Pin Counter**：记录正在访问这个page的Thread数量，thread访问该page前将counter + 1。如果一个page的pin counter大于0，则storage manager不允许将该page从buffer pool中剔除出去

数据库中的内存根据两种策略分配给buffer pool：

- **Global Policies**：DBMS为使整个工作负载收益而做出的决策，它会考虑所有活跃事务以找到分配内存的最佳决策
- **Local Policies**：它制定的决策将使单个查询或事务运行的更快。local policies将frame分配给特定事务，而不考虑并发事务的行为。

大多数系统结合两种策略。

## 2. Buffer Pool Optimizations

### 2.1 Multiple Buffer Pools

DBMS可以根据不同的目标维护不同的buffer pool，如对于page和index使用不同的buffer pool。这样每个buffer pool可以为其中存储的数据量身打造local policies。可以减少latch争用，提升locality。

获取record所在buffer pool的方法：

- **Object ID**：在record id中嵌入一个Object ID，然后建立Object ID到buffer pool的映射
- **Hasing**：对record id计算hash，再模buffer pool的数量得到在哪个buffer pool

### 2.2 Pre-Fetching

DBMS可以根据query plan进行page的预取。

假如query plan需要处理page1、page2、page3内的数据。query engine会先向buffer pool请求page1，buffer pool发现没有page1，于是阻塞线程执行磁盘IO，从disk中将page1读入。在query plan已知的情况下，buffer pool manager可以在page1数据执行时直接拷贝page2、page3到buffer pool中，之后query engine请求page2、page3时，buffer pool可以直接返回page地址而无需停顿下来进行磁盘IO。

一般在顺序访问page时使用。

### 2.3 Scan Sharing

查询可以重用存储计算和查询计算所产生的数据。允许多个查询连接到一个扫描table的游标。

如果一个查询开始扫描，而在开始之前已经有一个查询在做相同的事，那么DBMS会将第二个查询的游标附加到已经存在的查询游标上。DBMS会记录第二个查询的查询游标是在什么位置被附加到已存在的游标上的，这样可以帮助其完成整个查询操作。

### 2.4 Buffer Pool Bypass

即绕过buffer pool。

为避免过多的开销（如访问hash table），也为了避免污染现存的buffer pool（buffer pool中存储的数据在本次操作中可能不需要，但在未来将会很重要，也就是说当前查询的局部性会导致buffer pool的全局性受到较大影响），在执行顺序性扫描操作时，系统不会将fetched page存储在buffer pool中。相反，系统会为该查询单独分配一小块内存。如果操作需要读取磁盘上连续大量page，此方法效果很好。buffer pool bypass还可以用于临时数据（sort、join）。

## 3. OS Page Cache

大多数磁盘操作都通过OS API进行。除非另有明确说明，否则操作系统会维护自己的文件系统缓存。

大多数DBMS使用direct I/O `O_DIRECT` 绕过OS的cache，因为使用OS的cache会造成page的冗余副本，且页置换策略也不同。

## 4. Buffer Replacement Policies

置换策略用于决定从buffer pool中剔除哪一个frame（在buffer pool满时）。

### 4.1 LRU (LEAST_RECENTLY USED)

- 为每个page维护一个最近被使用的timestamp
- 需要剔除page时，选择具有最老timestamp的page
- 常规实现时map+bi-list，map由key映射value和指向bi-list的节点的指针，bi-list中存储key
  - 访问时，将bi-list中该节点删除，添加在bi-list头部，map中该key的指针指向头部
  - 剔除时，删除bi-list尾节点，将map中对应key删除

### 4.2 CLOCK

- LRU的一种近似设计，不需要为每个page都设置一个timestamp
- 每个page都有一个reference bit，当它被访问时，该bit置1
- 将page组织为有指针的环形buffer
  - 清除时，若page bit为1，置为0；否则剔除
  - 指针用来记录剔除的位置

### 4.3 Alternatives

LRU和CLOCK的问题：它们易遭受sequential flooding的影响，在顺序读取场景下，最近访问的page实际是最不需要的page。

解决sequential flooding问题的方法：

- **LRU-K**：跟踪每个page最近K个引用的历史记录作为 timestamp，以计算访问之间的间隔。此历史记录用于预测下一次将要访问page的时间。
- **Localization**：DBMS根据每个txn/query选择哪些page被剔除，这样可以最小化每次query对于buffer pool的污染。如Postgres对每个query维护一个private ring buffer。
- **Priority Hints**：允许事务在查询执行期间根据每个page的上下文告诉buffer pool该page是否重要