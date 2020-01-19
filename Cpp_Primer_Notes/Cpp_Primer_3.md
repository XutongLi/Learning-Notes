## 字符串、向量和数组

###  1. 作用域操作符（::）

编译器应从操作符左侧名字所示的作用域中寻找右侧的名字。

### 2. 头文件不应包含using声明

### 3. 初始化方式

**拷贝初始化** ：使用等号，编译器把等号右侧的初始值拷贝到新创建的对象中去（当初始值只有一个时）

**直接初始化** ：使用括号（当初始值有一个或多个时）

### 4. string size函数

size函数返回string对象的长度（即string对象中字符的个数）

size函数返回的是一个string::size_type类型的值，它是一个无符号类型的值且可以存放下任何string对象的大小。

可以通过`auto` 或者`decltype` 来推断变量的类型：

```C++
auto len = str.size();
```

如果一个表达式已经有了size函数就不要再使用int了，这样就可以避免混用int和unsigned可能带来的问题。

### 5. 字面值和string对象相加

string对象与字面值相加时，必须确保每个加法运算符两侧的运算对象至少有一个是string。

```C++
string s2 = s1 + "aaa"; 	//正确
string s3 = "hello " + "world";	//错误
string s4 = s1 + "hello " + "world";	//错误
```

字符串字面值与string是不同的类型。

### 6. cctype头文件

用于对字符进行判断或处理。

cctype头文件和ctype.h头文件内容是一样的，但是cctype中定义的名字从属于命名空间std。

### 7. vector可容纳类型

vector可容纳绝大多数类型的对象作为其元素，但是因为引用不是对象，所以不存在包含引用的vector。

### 8. 范围for循环与 vector添加元素

如果循环体内部包含有向vector对象添加元素的语句，则不能使用范围for循环。

（但凡是使用了迭代器的循环体，都不要向迭代器所属的容器添加元素）

### 9. 缓冲区溢出

通过下标访问不存在的元素

### 10. 迭代器

迭代器提供了对对象的间接访问，所有标准库容器都可以使用迭代器。

```C++
vector<int>::iterator it;	//it可以读写元素
vector<int>::const_iterator it2;	//it2只能读元素，不能写元素
vector<int> v;
auto a = v.begin();		//a表示v的第一个元素
auto b = v.end();		//b表示v尾元素下一个位置（尾后迭代器）
auto c = v.cbegin();	//c的类型是vector<int>::const_iterator
```

### 11. 数组不允许拷贝赋值

不能将数组的内容拷贝给其他数组作为其初始值，也不能用数组为其他数组赋值。

```c++
int a[] = {0, 1, 2};
int a2[] = a;	//错误
a2 = a;			//错误
```

### 12. 复杂的数组声明

```c++
int *ptrs[10];			//ptrs是含有十个指针的数组
int (*Parray)[10] = &arr;	//Parray指向一个含有十个整数的数组
int (&arrRef)[10] = arr;	//arrRef引用一个含有十个整数的数组
int *(&arry)[10] = ptrs;	//arry是数组的引用，该数组含有十个指针
```

### 13. size_t类型

使用数组下标的时候，通常将其定义为`size_t` 类型。size_t是一种机器相关的无符号类型，它被设计得足够大以便能表示内存中任意对象的大小。

定义在`cstddef` 头文件中。

### 14. string与char数组

#### char->string

- 允许使用以空字符结束的字符数组来初始化string对象或为string对象赋值
- string对象加法运算中允许使用以空字符结束的字符数组作为其中一个运算对象
- 在string对象的复合赋值运算（如+=）中，允许使用以空字符结束的字符数组作为右侧的运算对象

#### string->char

```c++
string s = "abc";
const char *str = s.c_str();
```

### 15. 数组初始化vector对象

```c++
int arr[] = {0, 1, 2, 3, 4, 5};
vector<int> vec(begin(arr), end(arr));	//begin和end计算arr的首指针和尾后指针
vector<int> subVec(arr + 1, arr + 4);	//拷贝arr[1]、arr[2]、arr[3]
```

### 16. 范围for循环处理多维数组

要使用范围for循环处理多维数组，除了最内层的循环外，其他所有循环的控制变量都应该是引用类型，这是为了避免数组被自动转成指针。

```c++
for (auto &row : ia)
    cor (auto col : row)	//要修改值时里层控制变量也是引用
```







































