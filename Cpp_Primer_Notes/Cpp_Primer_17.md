## 标准库特殊设施

### 1. tuple

当希望将一些数据组合成单一对象，又不想定义新的数据结构时，可以使用 `tuple` 。定义在 `tuple` 头文件中。

#### 定义和初始化

```c++
//tuple的构造函数时explicit的，必须用直接初始化语法
tuple<int, string, vector<int>> a(1, "aaa", {1, 2});
//make_tuple使用初始值的类型来推断tuple的类型
auto item = make_tuple("0-444-555", 1, 3.20);
```

#### 访问tuple的成员

`get<i>(t)` 函数返回第i个数据成员的引用。尖括号中的值必须是一个整型常量表达式。

```c++
auto book = get<0>(item);	//返回item第一个成员
get<2>(item) *= 0.8;		
```

获取tuple成员数量和类型：

```c++
typedef decltype(item) trans;	//trans是item的类型
//返回trans类型对象中成员的数量
size_t sz = tuple_size(trans)::value;	
//cnt的类型与item中第二个成员相同
tuple_element<1, trans>::type cnt = get<1>(item);	//也是从0开始计数
```

#### 关系和相等运算符

两个tuple有相同数量和类型的成员时可以比较。

由于tuple定义了<和==运算符，可以将tuple序列传递给算法，并且可以在无序容器中将tuple作为关键字类型。

### 2. bitset类型

标准库定义了 `bitset` 类，可以处理超过最长整型类型大小的位集合。在头文件 `bitset` 中。

#### 定义和初始化

bitset类有固定的大小：

```c++
bitset<32> bitvec(iU);	//32位，低位为1，其余为0
```

上述例子中，编号从0（从右开始）开始的二进制位被称为 **低位** ，编号到31结束的二进制位被称为 **高位** 。

#### 用unsigned值初始化bitset

当使用一个整型值来初始化bitset时，此值将被转换为 `unsigned long long` 类型并被当做位模式来处理。

bitset中的二进制位将是此模式的一个副本：

- 如果bitset的大小大于一个unsigned long long 的二进制位数，剩余的高位被置零
- 如果bitset的大小小于一个unsigned long long 中的二进制位数，则只使用给定值中的低位，超出bitset大小的高位被丢弃

```c++
//vec1比初始值小；初始值中的高位被丢弃
bitset<13> vec1(0xbeef);	//二进制序列为1111011101111
//vec2比初始值大，它的高位被置为0
bitset<20> vec2(0xbeef);	//二进制序列为00001011111011101111
//64位机器中，long long 0ULL是64个0， 所以~0ULL是64个1
bitset<128> vec3(~0ULL);	//0~63位为1,64~127为0
```

#### 用string初始化bitset

可以从一个string或者字符数组指针初始化bitset。字符串中下标最小的字符对应高位。

```c++
bitset<32> vec4("1100");
string str("11111110000000111");
bitset<32> vec5(str, 5, 4);		//从str[5]开始的4个二进制位
bitset<32> vec6(str, str.size() - 4);	//使用最后四个字符
```

如果string包含的字符数比bitset少，则bitset的高位被置为0。









