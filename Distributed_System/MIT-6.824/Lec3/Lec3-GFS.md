# Lec3 - GFS

**Big Storage System**

使用具备 **良好一致性** 的系统就像在与单个服务器交互

server 一次只执行一个 client 的请求，此情况下所有先前操作的数据反映了这些请求执行的顺序 。多以如果服务器收到一系列写请求操作，且以一定顺序去处理这些操作，当对请求读取时，会看到某种期望的 value 值 。

**容错**：构建主从复制服务器

人们对一致性的可接受范围取决于他们对那些异常行为的可接受度 。

**性能** ：将每一个文件数据由 **GFS** 自动拆分到很多服务器中，这样读写速度就会变得很快，当你从许多 client 读取一个文件时，你就会得到很高的吞吐量，这样也就能够读取比一个硬盘容量还大的单个文件了 。

**GFS** 被设计用于运行在单个数据中心的系统，所以并没有去讨论在全世界范围内放置数据副本 。

**GFS** 只对内服务，不对外销售 。适用于对大型数据（TB级或GB级）的连续访问而不是随机访问的场景 。

Google 在降低访问延迟上并没有做太多努力，他们的重心是放在了大吞吐量上面 。

**异端但非常有价值的观点**：让存储系统具备弱一致性是可行的 。（此论文中没有保证返回的数据是正确的）以此来获得更好的性能 。（如使用搜索引擎进行搜索，不会注意到返会的两万条结果中丢失了一条或展示顺序是错误的，这种系统对于错误数据的容忍度要远比银行之类的高得多）此外，GFS 会对数据进行校验和修正（checksum）。

**GFS structure**：

**master数据** 中，主要关心两张表，一张表管理了文件名和 chunk handles数组之间的映射关系，另一张表保存了每个chunk handle和chunk数据之间的映射关系 。这些都存在内存中。**nv** 表示非易失性 ，表示存在磁盘中 。

只有 primary 才有资格和 master 去进行过期时间判断 。master 会记住 lease 过期时间 （lease exprition）。

GFS 中，master 会将所有的操作记录以日志的形式放在磁盘中 ，并建立快照 。

**list of cs** 不需要存在磁盘，因为 master 重启后，会去和所有 chunksvr进行通信，并询问它们上面保存了哪些chunk 。

**primary** 不需要存在磁盘，因为 master 重启时就会忘记哪一个是primary，等待过期时间后，它可以为chunk指定一个新的 lease 。

这些意味着，当一个文件被追加了一个新的chunk，或者版本号变了（因为指定了新的primary），master必须现在它的日志里追加一条记录 。log使用文件而不是数据库的原因是：往文件尾追加效率很高，如果使用数据库，数据在磁盘上的分布是散的 。

master 时不时需要通过checkpoint将其完整状态写到磁盘 。master发生故障后重启时，要做的就是回滚到最近的checkpoint处，即回滚到该日志中创建这个checkpoint的时间点 。

**read过程**：

可能会多次读到同一chunk的不同连续区域，所以client将chunk handles和chunksvr list进行缓存 。

若 client 需求的字节范围跨越了一个chunk的边界，与client链接的库是一个GFS库，该库会直到如何将读请求分开，获取数据后再放到buffer中，返回给client 。

**write过程**：

client 包含一个文件名和一串字节，即想要写的数据包含在一个buffer中 。

master 会周期性地和chunkserver通信，询问这些服务器上持有哪些chunk以及它们的版本是是什么，若master记录的某chunk版本号为17，而它没找到含有此版本号的chunksvr，master不会响应，并告诉client端，”目前无法回答，请重试“ 。

使用 **版本号** 的理由是：master能根据它整理出哪些chunk服务器包含最新的chunk，master能授予这些chunk服务器称为primary的能力。

lease时间一般为60s，这其实是一种机制，它能确保我们最终不会有两个primary，防止一个primary挂掉，没有过期机制，在另一个上位primary之后，挂掉的那只重启后起冲突

只有当master指定了一个新的primary（即master不知道谁是primary时）时，版本号才会改变 。

防止 **split brain** （两个primary问题）：master提供了lease，即在特定时间里给chunk primary的权利，master知道该lease会持续多长时间 ，primary也直到lease会持续多长时间 。如果lease到期，则此chunk会直接拒绝Client请求，Client问master primary是谁，就会分配一个新的primary 。所以对于split brain问题，若master无法和primary通信，就会等到期时间过后，分配新的primary 。

**缺点**：

master并不是故障自动转移存储的，需要人工干预处理永久崩溃的master，可能要花费不少时间 。











