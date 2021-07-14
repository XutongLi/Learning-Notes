# Query Execution

## 1. Query Plan

![image](https://user-images.githubusercontent.com/29897667/124396950-5a253c80-dd3f-11eb-80ee-97d956e72617.png)

- DBMS将一个SQL语句转化成为一个**query plan**
- operators被组织在一个树中
- 数据流从叶节点流向跟节点
- 根节点的输出是query的结果
- 同一个query plan可以通过多种方式执行，尽可能使用index

## 2. Processing Model

一个DBMS的Processing Model定义了系统如何执行一个query plan。

### 2.1. Iterator Model

![image](https://user-images.githubusercontent.com/29897667/124756964-ac44a880-df5f-11eb-85c5-4b8b570c9c25.png)

![image](https://user-images.githubusercontent.com/29897667/124756833-7f909100-df5f-11eb-81cb-b0999cc2ac19.png)

也叫作 **Volcano Model** 或 **Pipeline Model**。

- 叫pipeline原因：对于一个tuple，这种模型能使它在query plan中尽可能多地被使用，即在一个operator中处理完后，然后返回并传入下一个operator继续处理（让一个tuple经过尽可能多的operators）。对一个tuple进行一些列处理加工的过程称为pipeline
- 每个query plan operator实现一个`NEXT` 函数
  - 每次调用 `NEXT` ，operator会返回一个tuple或null（没有tuple时）
  - operator会实现一个loop调用子节点的`NEXT`方法获取tuple并处理它们
- 一些操作会阻塞pipeline直到children发出它的所有tuples，这些操作称为**pipeline breaker**，包括：`join`、 `subqueries` 、`order by`
- output control，如limit很容易做，因为在获取到需要的tuples后就可以停止调用`NEXT`

### 2.2. Materialization Model

![image](https://user-images.githubusercontent.com/29897667/124759116-e020cd80-df61-11eb-866e-e6eeb61a1991.png)

- 主要用于内存型数据库
- 每个operator一次处理从子节点的`OUTPUT`方法获取的所有tuples
- operator返回它要发出的所有tuples
  - operator `materialize(物化)` 它的输出称为单个结果，可以是一个materialized row或单个column
  - 在operator执行完成之前，不会返回及获取更多数据
  - DBMS可以下推`hints(提示)`（沿着query plan tree）取避免扫描太多tuples
- 这种方法对OLTP workload比较好，因为query仅一次访问一小部分tuples；对于会产生大量中间结果的OLAP查询不太好，因为DBMS可能需要将这些中间结果溢出到disk。

### 2.3. Vectorized / Batch Model

![image](https://user-images.githubusercontent.com/29897667/124763192-49a2db00-df66-11eb-8150-07f09e5dbf4f.png)

- 也实现`NEXT`方法，但是发出一批tuples而非单个tuple
- operator实现会优化为对一批数据的处理（operator的内部循环一次处理多个tuples）
- 主要用于分析型数据库，因为OLAP一次需要扫描大量的tuples，使用vectorized model可以调用更少次数的NEXT
- 现代CPU有SIMD指令，允许一次对一堆数据进行多个操作。如果有一堆数据，这堆数据可以放在CPU寄存器当中，通过这种单条指令，可以高效地对这些数据进行条件判断或计算。

### 2.4. Plan Processing Direction

#### 2.4.1. Top-to-Bottom

- 从根节点开始，从子节点获取数据
- tuples通过调用函数来获取

#### 2.4.2. Bottom-to-Top

- 从叶节点开始，将数据push给父节点

- 在向上传递数据时，要确保正在处理的数据能够放在CPU缓存和寄存器中
- 只能用于in-memory数据库

## 3. Access Methods

**Access Method** 指DBMS从table中访问数据的方式，即在query plan tree中叶子节点发生的事。

### 3.1. Sequential Scan

#### 3.1.1. 基本方法

![image](https://user-images.githubusercontent.com/29897667/124976624-1a719400-e062-11eb-83dc-a9115cda0276.png)

- 对table中的每个page，从buffer pool中获取它，获取其中的每个tuple并进行后续操作
- DBMS维护内部的 **cursor** 来指向最后访问的page或slot

#### 3.1.2. 优化

**Prefetching**：预先取当前page后连续的page，防止DBMS因为从disk上取page而造成的阻塞

**Buffer Pool Bypass**：使用一个小buffer对查询进行缓存，避免污染buffer缓存（避免sequential flooding问题）

**Parallelization**：使用多线程或多进程并行地执行scan

**Zone Map**：

![image](https://user-images.githubusercontent.com/29897667/125164277-f62dc880-e1c3-11eb-9c93-b347152702c3.png)

- 对于每个page的属性值计算一些聚合信息（如MIN、MAX、AVG、SUM、COUNT等）
- DBMS通过检查zone map来确定是否要循序访问这个page
- 有的系统将zone map保存在当前page中，有的保存在专门的zone map page中，上面保存了不同page的zone map
- 更新page中的tuple时就会更新zone map
- 一般用在OLAP中

**Late Materialization**：

- 延迟将数据从一个operator传播到下一个operator
- 在列存中，单个tuple的各个属性不能被一次获取，且有的属性值在plan tree中的一个节点中用到，因此先在operator中传递record_id，延迟对整个tuple的读

**Heap Clustering**：

- tuples顺序保存在索引的叶节点中
- 当查询使用到聚集索引的属性时，DBMS可以直接跳转到需要访问的page上

****

### 3.2. Index Scan

#### 3.2.1. Basic

**基本思路**：DBMS选择合适的索引取寻找query所需要的tuples

**index的选择取决于**：

- index包含哪些属性
- query查找哪些属性
- query的判断条件（属性值范围）
- index是否是unique index

#### 3.2.2. Multi-Index Scan

通过不同的索引进行多次查找，然后基于判断条件，将结果进行合并（AND、OR）

如果对于一次query可以使用多个index

- 计算使用每个index的record id set（set可以是bitmap、hash table或Bloom Filter等形式）
- 基于query的预测组合这些sets（union or intersect）
- 最后检索records

![1626026698300](C:\Users\XutongLi\AppData\Roaming\Typora\typora-user-images\1626026698300.png)

按照非聚集索引中的顺序查找tuples会导致随机IO，因此DBMS可以首先查找出所有tuples并且基于page id排序（找出tuple对应的主键key，并将逐渐key排序，在进行disk IO，即 [MySQL MRR](https://github.com/XutongLi/My_Document/issues/9#issuecomment-874611801)）

## 4. Expression Evaluation

DBMS将 **WHERE Clause** 表示为 **expression tree**。（在一个operator中）

表达树上不同节点代表了条件判断中不同类型的表达式。

![image](https://user-images.githubusercontent.com/29897667/125205672-01f9b780-e2b6-11eb-9edc-053d2e241596.png)

为了在运行时评测expression tree（的条件判断的正确性），DBMS会维护一个context handler去包含执行的一些元数据（current tuple、table schema、parameter）。这些DBMS会便利tree去评测操作并产生结果。

## 5. Parallel Execution 的好处

- 想在数据库中利用CPU或GPU的多核
- 更高的吞吐量：每秒执行更多的查询、每秒处理更多的数据
- 更小的延迟：单个查询要花的执行时间变小
- 更好的系统响应能力和可用性：系统对请求更快地响应

## 6. 并行(parallel)和分布式(distributed)

### 6.1. 相同点

将一个数据库分散到多个资源（多台机器、多个CPU、多个磁盘等）上，以改善性能、成本、延迟等。对外表现为单个DB实例。

### 6.2. 不同点

**Parallel DBMS**

- 资源在物理上相邻
- 资源之间的通信方式高速、便宜、可靠

**Distributed DBMS**

- 资源在物理上距离比较远 
- 资源之间的通信方式慢、且成本和可靠性不可忽视

## 7. Parallel Process Model

如何组织系统来通过多个worker处理并发请求。情况包括一个请求分散给多个worker、多个请求给多个worker。

### 7.1. Process per Worker

- 每个worker是一个独立的OS进程，它依赖于OS sheduler
- 使用shared memory来存储全局数据结构
- 一个进程crash不会影响到整个系统

### 7.2. Process Pool

- 一个worker使用pool中空闲的worker，不会为进来的每个连接去创建一个进程
- 仍然依赖于OS scheduler和shared memory
- 这种方法对CPU缓存一致性不好，因为不能保证在查询间使用一个进程

### 7.3. Thread per Worker

- 单个进程中有多个worker threads
- DBMS使用一个dispatcher thread进行调度
- 一个thread crash会导致系统crash
- 优点：上下文切换开销更小，且不需要管理shared memory

## 8. Parallel类型

### 8.1. Inter-Query

- 同时执行不同的queries
- 可以提升吞吐率并减少延迟
- 如果queries是只读的，则容易做到；如果queries会同时更新DB，则很难正确做到（并发控制）

### 8.2. Intra-Query

- 将一个查询拆分为多个子任务或片段，然后在不同的worker上同时并行执行这些任务
- 可以减少长时间运行的query的延迟，对于分析性查询优化效果明显
- 每个operator既接受数据也生产数据，可以将其看作“生产者-消费者模型”
- 对于每个relational operator都有并行算法：
  - 可以使用多个thread访问中心数据结
  - 或者将收到的数据进行划分，然后使用多线程分别处理这些数据分区（这样就不需要对worker进行协调了）
- 如 **Parallel Grace Hash Join**，可以让每个worker负责一个bucket中的match

![image](https://user-images.githubusercontent.com/29897667/125506043-dc7562ac-7cda-462d-bbee-1c3393a19a44.png)

#### 8.2.1. Intra-Operator Parallelism (Horizontal)

- 将一个完整的操作拆分成多个并行的操作，即将操作的数据分为多段，每一段的执行函数都是一样的，每段中的数据为输入数据中的一部分
- 如对B+Tree的parallel scan
- 使用叫**exchange**的operator将这些结果组合在一起，它们放在query plan中可以并行执行的位置
- **exchange**阻止DBMS在plan中执行它上面的操作符，直到它从子级接收到所有的数据

**exchange operators**：

- **Gather**：将多个workers的结果组合成一个输出流（多对一）
- **Repartition**：将多个输入流重新组织成多个输出流，即以一种方式partition，再以另一种方式redistribute（多对多）
- **Distribute**：把一个单一输入流分成多个输出流

**示例**（图中维护的为一个hash table）：

![image](https://user-images.githubusercontent.com/29897667/125512826-d418c1e0-c89b-4adb-b1fa-0744a2aef2ec.png)

#### 8.2.2. Inter-Operator Parallelism (Vertical)

- 不同的线程在同一时间执行不同的operator
- 将数据从一个阶段传输到下一个阶段而不进行物化（生成临时表），又称为**pipelined parallelism**

**示例**：

![image](https://user-images.githubusercontent.com/29897667/125514710-786d11d8-8a1a-49d9-b0ab-256b7a00b06f.png)

#### 8.2.3. Bushy Parallelism

- inter-operator parallelism的扩展，workers同时执行query plan的不同分段的多个operators
- 依然使用exchange operator从各段组合中间结果
- 每个operator就是它自己的worker

**示例**：

![image](https://user-images.githubusercontent.com/29897667/125516220-c2ae9aff-c25e-46eb-ac83-7c928258128d.png)

![image](https://user-images.githubusercontent.com/29897667/125516179-eadaf2fc-dc41-4f34-9e74-0bb03d9d7acd.png)

## 9. IO Parallelism

如果disk IO是瓶颈的话，那么并行执行查询也不会带来太大性能提升。

对数据库文件和数据进行拆分，分散到存储设备的不同位置处（让多个存储设备以单个逻辑设备的形式来供数据库系统使用）

**划分粒度**：

- 每个DB多个disk
- 一个disk上一个DB
- 一个disk上一个relation（table）
- 划分relation到多个disk上

### 9.1. Multi-Disk Parallelism

- 使用多块存储介质来存数据库文件
- 可通过RAID配置实现，对于DBMS是透明的
- 不可以使用多worker并行访问，因为DBMS不知道底层存储布局

### 9.2. File-based Parallelism

- 可对每个database指定disk上的位置
- buffer pool manager将page映射到disk位置上

### 9.3. Logical Partitioning

- 将单个logical table划分为不相交的physical partition，这些Partition都被单独管理。
- 这些Partitioning对于应用是透明的，即应用访问这个logical table无需关注它的具体存储方式

#### 9.3.1. Vertical Partitioning

- 分开存储一张表上的属性（类似于列存）
- 需要取存储tuple信息以重建原始record

![image](https://user-images.githubusercontent.com/29897667/125622283-f84d9c22-6f54-4d76-a134-2bf396bc5d10.png)

#### 9.3.2. Horizontal Partitioning

- 基于一些Partitioning keys，将一个table的tuples划分为不相交的分段
- 有不同的方法来决定如何做Partition
  - Hash Partitioning
  - Range Partitioning
  - Predicate Partitioning
- 划分方式的有效性取决于query

![image](https://user-images.githubusercontent.com/29897667/125622315-463dec54-1113-488a-90fc-c0b63a77954a.png)

