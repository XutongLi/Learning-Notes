# Hash Table

## 1. 概念
**Hash Table** 是一个无序的关联数组，用于key到value的映射。通过 **Hash Function (h)** 可以计算出Key在Hash Table中的的offset，从而找到value。

空间复杂度为 `O(n)` ，时间复杂度为 `avg O(1)`、`worst O(n)`。

**桶 (Bucket)** 表示能存储一条或多条记录的一个存储单位。通常一个桶就是一个磁盘块。

Hash可用于：
- **散列文件组织 (hash file organization)**：通过对key计算hash function直接获得包含该记录的磁盘块地址。
- **散列索引组织 (hash index organizaiton)**：把key以及它相关联的指针组织成一个散列文件结构。

## 2. Hash Function
对于任意key，hash function返回该key的一个整型表示。通常hash function在key中字符的内部二进制机器表示上执行计算。

hash function 应该尽可能快以及有尽可能小的碰撞率，无加密需求。hash function应该满足两个特性：
- 均匀的：hash function应该从所有可能的key集合中为每个桶分配相同数量的key。
- 随机的：hash function不应与key的任何外部可见的排序相关，即无论key值如何分布，每个桶分布到的key数目几乎相同。

对于10字节以下的key，Google CityHash和FarmHash比较快；对于更长的key，Facebook的XXHash比较快。

## 3. Static Hash Table

### 3.1 概念
关联数组的大小是固定的。为方便表述，令 *nr* 为要存储的记录总数，*f* 为一个桶中可以存储的记录数量，*nb* 为桶的数量。

当出现桶不足或者偏斜（某些桶分配到的记录比其他桶多）情况时，就会出现 **桶溢出** 情况。为减少桶溢出情况，一般使 `nb = (nr/fr)*(1+d)` ，*d* 为避让因子，典型值为0.2。

**解决桶溢出问题** 的方法有：
- **闭地址法**：链地址法（溢出桶）
- **开地址法**：当一个桶满了后，系统将记录插入到初始桶集合的其他桶中，方法有线性探测法、再哈希法等。

闭地址法适合用于数据库，开地址法因为删除麻烦，一般用于构造编译器和汇编器的符号表。

对于有 **重复key** 的情况，可以使用：
- 独立链表法：对每个key，将其所有value存在一个独立存储区域中
- 冗余key法：将重复key对应的entry一起存在hash table中

### 3.2 解决桶溢出的特殊方法 (*fr=1*)

#### Robin Hood Hashing
线型探测法的一种变体，尽量平均key实际存储offset与其最优offset之间的距离。
- 对于每个key都跟踪其在hash table中的实际offset以及最优offset之间的距离dis
- 当key A的dis>key B的dis时，key A抢占key B的offset

#### Cuckoo Hashing
使用含有不同hash function seed的hash table。
- 插入时，检查所有的hash table，并找到一个有空位的插入
- 如果所有hash table都没有针对此key的空位，剔除其中一个，将新key插入，就key重新hash寻找空位
- 删除和查询操作都是O(1)，因为每个hash table只有一个空位被查找

### 3.3 缺点
static hash table要求DBMS预先知道数组大小，否则按需扩大或缩小数组。
解决方法：Dynamic Hash Table

## 4. Dynamic Hash Table

### 4.1 Extendible Hashing

可扩充散列可以通过桶的分裂或合并来适应数据库大小的变化，这样可以保持空间的使用效率。由于重组每次仅作用于一个桶，因此所带来的性能开销较低。

在使用extendible hash时，hash function *h* 将key映射为一个 *b* 位整数，一个典型的 *b* 值为32。

对于桶地址表（数组），有一个散列前缀 *i (0<=i<=b)*，一开始不使用散列值的所有位数，*i* 值随着数据库大小的变化而增大或者减少。

因为有共同散列前缀的几个表项可能指向同一个桶，因此每个桶也有一个前缀值 *ij*，表明确定该桶需要的散列值位数。

具体操作：
- **确定key对应桶的位置**：取 *h(key)* 的前 *i* 个高位，这个数为数组offset，再从数组表项中得到指向桶 *j* 的指针。
- **插入**：若该桶未满，直接插入；否则分裂该桶，将该桶中的entry重新分配：
    - 若 *i = ij*：此时指向该桶的只有一个表项，需要将地址表大小扩大一倍，*i+=1*，原表中的每个表项都由两个表项取代，两个表项都含有和原始表项一样的指针。现在有两个表项指向 *j*，系统分配新桶 *z*，令第二个表项指向此新桶，令 *ij* 和 *iz* 都为 *i*，将桶 *j* 中的记录重新散列。分裂后再尝试插入，若插入不成功则继续分裂。若桶 *j* 中所有key一样，此时继续插入已无作用，需使用溢出桶方法。
    - 若 *i > ij*，则系统不需要扩大一倍，直接分裂。
- **删除**：确定好 key 对应的桶 *j* 后，把记录从桶中删除。若桶变空，则桶删除，此时有的桶可能需要合并，地址表大小需要减半。只有当桶数目减少很多时，减小桶地址表的大小才是值得的（改变表大小开销很大）。


### 4.2 Linear Hashing

https://blog.csdn.net/jackydai987/article/details/6673063

hash table维护一个指针，该指针指向下一个要被分裂的桶。有桶满时并非直接在该桶上分裂，而是在指针指向的桶上分裂（这么做是为了方便判断哪些桶被分裂了）
- 当一个桶满时，在指针指向的位置分裂桶。并创建一个新的hash function，将该桶中key按照新的hash function重新分配。（此时若分裂的不是满的桶，则新插入的元素先存在溢出桶中，直到该桶被分裂）
- 查找新key位置时，如果hash function映射到分裂指针指向的桶前，则应用新的hash function重新计算。
- 当分裂指针到达最后一个桶时，将分裂指针置为0，删除旧的hash function并用新的hash function取代它

Linear Hashing 具有无限扩张能力，支持O(1)的查找。