## 关联容器

### 1. 关联容器与顺序容器区别

关联容器中的元素是按关键字来保存和访问的；顺序容器中的元素按它们在容器中的位置来顺序保存和访问。

### 2. 关联容器类型

| 类型               | 描述                             | 是否有序 | 头文件        |
| ------------------ | -------------------------------- | -------- | ------------- |
| map                | 关联数组，保存关键字-值对        | 有       | map           |
| set                | 关键字即值，即只保存关键字的容器 | 有       | set           |
| multimap           | 关键字可重复出现的map            | 有       | map           |
| multiset           | 关键字可重复出现的set            | 有       | set           |
| unordered_map      | 用hash函数组织的map              | 无       | unordered_map |
| unordered_set      | 用hash函数组织的set              | 无       | unordered_set |
| unordered_multimap | hash组织的map，关键字可重复出现  | 无       | unordered_map |
| unordered_multiset | hash阻止的set，关键字可重复出现  | 无       | unordered_set |

### 3. 有序容器关键字类型要求

有序容器关键字类型必须定义元素比较的方法。默认情况下，标准库市容关键字类型的 `<` 运算符来比较两个关键字。

实际编程中，如果一个类型定义了 `<` 运算符，则它可以用作关键字类型。

用尖括号指出要定义哪种类型的容器，自定义的操作类型必须在尖括号中紧跟着元素类型给出。比较操作类型是一种函数指针类型。

```c++
bool compareISBN(const Sales_data &lhs, const Sales_data &rhs) 
    return lhs.isbn() < rhs.isbn();
//bookstore中多条记录可以拥有相同的ISBN
//bookstore中元素以ISBN顺序进行排序
multiset<Sales_data, decltype(compareISBN)*> bookstore(compareISBN);
```

用compareISBN来初始化bookstore对象，这表示向bookstore添加元素时，通过调用compareISBN来为这些元素排序。

### 4. pair类型

标准库类型 `pair` 是map的基本类型，头文件 `utility` 。有两个public的数据成员，分别命名为 `first` 和 `second` ，用普通的成员访问符号访问它们。

`make_pair` 可用来生成pair对象。

  ### 5. 关联容器类型别名

`key_type` 是容器的关键字类型

`mapped_type` 是每个关联字关联的类型，只适用于map

`value_type` 对于set，与key_type相同；对于map，为 `pair<const key_type, mapped_type>`

### 6. 关联容器添加元素

#### 不含重复关键字容器

对于 **不包含重复关键字** 的容器，添加 **单一元素** 的 `insert` 和 `emplace` 返回一个 `pair` 。
`pair` 的 `first` 成员是一个迭代器，指向具有给定关键字的元素； `second`  成员是一个bool值，指出元素插入成功还是已经存在于容器。

如果关键字已在容器中，则 `insert` 什么也不做，返回一个指向该关键字元素的迭代器和false。如果关键字不存在，元素被插入容器中，且bool值为true。

#### 包含重复关键字容器

对于 **允许重复关键字的容器** ，接受 **单个元素** 的 `insert` 操作返回一个指向新元素的迭代器。

### 7. map下标操作

`c[k]` 返回关键字为k的元素。如果k不在c中，添加一个关键字为k的元素，对其值进行初始化。

`c.at(k)` 访问关键字为k的元素，带参数检查；若k不在c中，抛出一个out_of_range异常。

`multimap` 和 `unordered_map` 没有下标操作，因为这些容器可能有多个值和一个关键字相关联。

### 8. 关联容器访问元素

对于不允许重复关键字的容器，使用 `find` 和 `count` 没什么区别。

#### 允许关键字重复容器查找

如果一个 `multimap` 和 `multiset` 中有多个元素具有给定关键字，则这些元素在容器中会相邻存储。

**方法一**：find + count

```c++
string item("Brian");
auto cnt = authors.count(item);		//元素数量
auto it = authors.find(item);		//第一本书
while (cnt) {
    cout << it->second << endl;
    ++ it;	--cnt;
}
```

**方法二** ：lower_bound + upper_bound

`lower_bound` 返回的迭代器将指向第一个具有给定关键字的元素，而 `upper_bound` 返回的迭代器则指向最后一个匹配给定关键字的元素之后的位置。

如果没有元素与给定的关键字匹配，两者会返回相等的迭代器——都指向给定关键字的插入点，能保证容器中元素顺序的插入位置。

```c++
for (auto beg = authors.lower_bound(search_item), end = authors.upper_bound();
     beg != end; ++ beg) 
    cout << beg->second << endl;
```

**方法三** ：equal_range

`equal_range` 函数接受一个关键字，返回一个迭代器pair。若关键字存在，则第一个迭代器指向第一个与关键字匹配的元素，第二个迭代器指向最后一个匹配元素之后的位置。若未找到匹配元素，则两个迭代器都指向关键字可以插入的位置。

```c++
for (auto pos = authors.equal_range(item); pos.first != pos.second; ++ pos.first)
    cout << pos.first->second << endl;
```

































