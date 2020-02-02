## 重载运算与类型转换

### 1. 运算符重载

#### 定义

当运算符作用于类类型的对象时，可以通过运算符重载重新定义该运算符的含义。

重载的运算符是具有特殊名字的函数，它们的名字由关键字 `operator` 和运算符号构成。重载的运算符也包含返回类型、参数列表以及函数体。

#### 参数数量

重载运算符函数的参数数量与运算符作用的运算对象数量一样多。

如果一个运算符函数是成员函数，则它的第一个运算对象绑定到隐式的this指针上。成员运算符函数的显式参数数量比运算对象数量少一个。

#### 什么样的运算符可以被重载

对于一个运算符函数来说，它或是类的成员，或者至少好友一个类类型的参数。

只能重载已有的运算符。

通常情况下，**不应该** 重载逗号、取地址、逻辑与和逻辑或运算符。

#### 直接调用重载的运算符函数

```c++
//一个非成员运算符函数的等价调用
data1 + data2;
operator+(data1, data2);
//一个成员运算符函数的等价调用
data1 += data2;
data1.operator+=(data2);
```

### 2. 输入和输出运算符

#### 重载输出运算符 <<

输出运算符的 **第一个形参** 是 **非常量ostream对象的引用** 。（非常量是因为向流写入内容会改变其状态；是引用是因为无法直接复制一个ostream对象）

**第二个形参** 是一个 **常量的引用** 。（是引用的原因是避免复制形参；常量是因为打印对象不会改变其内容）

```c++
ostream &operator<<(ostream &os, const Sales_data &item) {
    os << item.isbn();
    return os;
}
```

运算符返回刚使用的ostream的引用。

输出运算符负责打印对象的内容而非控制格式。

输入输出运算符必须是非成员函数，因为IO运算符常需要读写类的非公有数据成员，所以IO运算符一般被声明为友元。

#### 重载输入运算符 >>

输入运算符的 **第一个形参** 是运算符将要读取的流的引用。**第二个形参** 是将要读入到的（非常量）对象的引用。

该运算符返回给定输入流的引用。

```c++
istream &operator>>(istream &is, Sales_data &item) {
    double price;
    is >> item.isbn >> item.sold >> price;
    if (is)	//检查输入是否成功
        item.rev = item.sold * price;
    else
        item = Sales_data();	//输入失败，对象被赋予默认的状态
    return is;
}
```

输入运算符必须处理输入可能失败的情况，而输出运算符不需要。

### 3. 算术和关系运算符

算术和关系运算符一般定义为非成员函数以允许对左右侧的运算对象进行转换。形参为常量引用。

#### 算术运算符

如果类同时定义了算术运算符和相关的复合赋值运算符，通常情况下使用复合赋值来实现算术运算符：

```c++
//复合赋值运算符（成员函数）
Sales_data& Sales_data::operator+=(const Sales_data &rhs) {
    units_sold += rhs.units_sold;
    revenue += rhs.revenue;
    return *this;
}
//算术运算符函数（非成员函数）
Sales_data operator+(const Sales_data &lhs, const Sales_data &rhs) {
    Sales_data sum = lhs;
    sum += rhs;
    return sum;	
}
```

#### 相等运算符

如果一个类定义了 `operator==` ，则这个类也应该定义 `operator!=` 。

```c++
bool operator==(const Sales_data &lhs, const Sales_data &rhs) {
    return lhs.isbn() == rhs.isbn() && lhs.revenue == rhs.revenue;
}
bool operator!=(const Sales_data &lhs, const Sales_data &rhs) {
    return !(lhs == rhs);
}
```

#### 关系运算符

如果存在唯一一种逻辑可靠的 `<` 定义，则应该考虑为这个类定义 `<` 运算符。如果类同时还包含 `==` ，则当且仅当 `<` 的定义和 `==` 产生的结果一致时才定义 `<` 运算符。

#### 赋值运算符

赋值运算符必须定义为成员函数，返回左侧运算对象的引用。

除拷贝赋值运算符和移动赋值运算符外，还可以定义其他赋值运算符以使用别的类型作为右侧运算对象。

```c++
//花括号赋值
StrVec &StrVec::operator=(initializer_list<string> il) {
    auto data = alloc_n_copy(il.begin(), il.end());
    free();		//同样要先释放自己的空间
    elements = data.first;
    first_free = cap = data.second;
    return *this;
}
StrVec v = {"aaa", "bbb"};
```

### 4. 下标运算符

表示容器的类一般定义下标运算符 `operator[]` ，用以通过位置访问元素。

下标运算符必须是成员函数。

```c++
class StrVec {
public:
    string& operator[](size_t n)		//返回引用，可修改
        return elements[n];
    const string& operator[](size_t n) const	//返回常量引用，不可修改
        return elements[n];
};
const StrVec cevc = svec;	//假设svec是一个StrVec对象
svec[0] = "aaa";	//正确
cvec[0] = "aaa";	//错误
```

### 5. 递增递减运算符

定义递增递减运算符的类应该同时定义前置版本和后置版本。这些运算符通常应被定义为类的成员。

#### 前置递增/递减运算符

前置运算符返回递增或递减后对象的引用：

```c++
class StrBlobPtr {
public:
    StrBlobPtr& operator++();
    StrBlobPtr& operator--();
};
StrBlobPtr& StrBlobPtr::operator++() {
    check(curr, "error");	//若curr已到达容器尾后元素，则无法递增它
    ++ curr;
    return *this;
}
StrBlobPtr& StrBlobPtr::operator--() {
    -- curr;	//若curr为0，继续递减产生一个无效下标
    check(curr, "error");	//检查
    return *this;
}
```

#### 后置递增/递减运算符

后置版本接收一个额外的（不被使用）的int类型的形参作为与前置版本的区分。编译器为这个形参提供一个值为0的形参。

后置运算符应返回对象的原值（递增或递减之前的值），返回的形式是一个值而非引用。

```c++
class StrBlobPtr {
public:
    StrBlobPtr operator++(int);
    STrBlobPtr operator--(int);
};
StrBlobPtr StrBlobPtr::operator++(int) {
    StrBlobPtr ret = *this;	//记录之前的状态
    ++ *this;	//调用前置运算符，检查在其中操作
    return ret;	//返回之前的状态
}
StrBlobPtr StrBlobPtr::operator--(int) {
    StrBlobPtr ret = *this;
    -- *this;	
    return ret;
}
```

### 6. 成员访问运算符

箭头运算符必须是类的成员。解引用运算符通常是类的成员。

```c++
class StrBlobPtr {
public:
    string& operator*() const {
        auto p = check(curr, "error");
        return (*p)[curr];	//*p是对象所指的vector
    }
    string* operator->() const {
        return & this->operator*();	//调用解引用运算符并返回解引用结果元素的地址
    }
};
```

将这两个运算符定义成const成员，因为获取一个元素不会改变对象的状态。

### 7. 函数调用运算符

类重载了函数调用运算符，可以像使用函数一样使用该类的对象。

函数调用运算符必须是成员函数。一个类可以定义多个不同版本的调用运算符，相互之间应该在参数数量或类型上有所区别。

```c++
struct absInt {
    int operator()(int val) cosnt {
        return val < 0 ? -val : val;		//返回绝对值
    }    
};
int i = -42;
absInt obj;		//obj称为函数对象
int ui = obj(42);	//返回42
```

#### lambda是函数对象

编写一个lambda表达式之后，编译器将该表达式翻译成一个未命名类的未命名对象。lambda表达式产生的类中含有一个重载的函数调用运算符。

```c++
stable_sort(words.begin(), words.end(),
           [](const string &a, const string &b) {return a.size() < b.size();});
//等价于：
class ShorterString {
public:
    bool operator() (const string &s1, const string &s2) const {
        return s1.size() < s2.size();
    }
};
stable_sort(words.begin(), words.end(), ShorterString());
```

当一个lambda表达式 **通过引用捕获变量** 时，编译器可直接使用该引用而无需再lambda产生的类中将其存储为数据成员。

**通过值捕获的变量** 被拷贝到lambda中，因此这种lambda产生的类必须为每个值捕获的变量建立对应的数据成员，同时创建构造函数，令其使用捕获的变量的值来初始化数据成员。

```c++
auto wc = find_if(words.begin(), words.end(),
                 [sz](const string &a) {return a.size() >= sz;});
//等价于：
class SizeComp {
public:
    SizeComp(size_t n) : sz(n) {}
    bool operator() (const string &s) const {
        return s.size() >= sz;
    }
private:
    size_t sz;
};
auto wc = find_if(words.begin(), words.end(), SizeComp(sz));
```





 











