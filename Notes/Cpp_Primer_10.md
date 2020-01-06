## 泛型算法

### 1. 泛型算法实际操作

泛型算法并不直接操作容器，它们只会运行于迭代器之上，执行迭代器的操作。算法可能改变容器中保存的元素的值，也可能在容器内移动元素，但永远不会直接添加或删除元素。

### 2. 迭代器参数

一些算法从两个序列中读取元素。构成着两个序列的元素可以来自于不同类型的容器。而且，两个序列中的元素的类型也不要求严格匹配。算法要求的只是能够比较两个序列中的元素。

### 3. 谓词

谓词是一个可调用的表达式，其返回结果是一个能用作条件的值。标准库算法所使用的谓词分为两类：一元谓词（只接受单一参数）和二元谓词（接受两个参数）。接受谓词参数的算法对输入序列中的元素调用谓词。

```c++
bool isShorter(const string &s1, const string &s2) {
    return s1.size() < s2.size();
}
sort(words.begin(), words.end(), isShorter);	//按照长度由短到长排序words
```

### 4. stable_sort

稳定排序算法，维持相等元素的原有顺序。

### 5. lambda表达式

#### 概念

当希望进行的操作需要更多参数，超出了算法对谓词的限制时，可以使用lambda表达式。

#### 形式

一个lambda表达式表示一个可调用的代码单元，具有一个捕获列表、一个返回类型、一个参数列表和一个函数体。它可以定义在函数内部。

```c++
[capture list] (parameter list) -> return type {function body}
auto f = [] {return 42;}
```

`capture list` （捕获列表）是一个lambda所在函数中定义的局部变量的列表，如果为空表示此lambda不使用它所在函数中的任何局部变量。

捕获列表只用于局部非static变量，lambda可以直接使用局部static变量和它所在函数之外声明的名字。

捕获也分为值捕获（拷贝发生在lambda创建时，而不是调用时）和引用捕获。可以使用隐式捕获，让编译器根据lambda体中的代码来推断我们要使用哪些变量。

`return type` 必须使用尾置返回。可以省略返回类型，此时lambda根据函数体中的代码推断出返回类型。

`parameter list` 可以省略。不可以有默认参数，所以lambda调用的实参数目永远与形参数目相等。

`function body` 与普通函数一样，若不包含return且未指定返回类型，返回void。

#### 调用

```c++
cout << f() << endl;	
```

lambda的调用方式与普通函数的调用方式相同，都是使用调用运算符。

#### 例子

```c++
//打印words中长度大于sz的string
void biggies(vector<string> &words, vector<string>::size_type sz) {
    elimDups(words);	//将words按字典序排序，且删除重复单词
    //按长度排序，长度相同的单词维持字典序
    stable_sort(words.begin(), words.end(), 
                [](const string &a, const string &b) {return a.size() < b.size();});
    //获取一个迭代器，指向第一个满足长度>=sz的元素
    auto wc = find_if(words.begin(), words.end(), 
                      [sz](const string &a) {return a.size() >= sz;});
    //计算满足size >= sz的元素数目
    auto cnt = words.end() - wc;
    for _each(wc, words.end(),
             [](cosnt string &s) {cout << s << " ";});
    cout << endl;
}

//错误，编译器推断lambda返回类型为void，而实际上返回了int
transform(v.begin(), v.end(), v.begin(), 
          [](int i) {if (i < 0) return -i; else return i;});	
//以下两种写法正确
transform(v.begin(), v.end(), v.begin(),
         [](int i) {return i < 0 ? -i : i; });
transform(v.begin(), v.end(), v.begin(),
         [](int i) -> int {if (i < 0) return -i; else return i;});//定义尾置返回类型
```

### 6. 参数绑定

#### 概念

`bind` 函数是一个通用的函数适配器，它接受一个可调用对象，生成一个新的可调用对象来”适应“原对象的参数列表。它是一个标准库函数，定义在头文件`functional` 中。

#### 形式

```c++
auto newCallable = bind(callable, arg_list);
```

`newCallable` 本身是一个可调用对象，`arg_list` 是一个逗号分隔的参数列表，对应给定的`callable` 的参数。即，当我们调用 newCallable 时，newCallable会调用Callable，并传递给它arg_list中的参数。

arg_list中可包含占位符（_n），\_1为newCallable的第一个参数。

#### 举例

```c++
//check6是一个可调用对象，接受一个string类型的参数，并用此string和值6来调用check_size
auto check6 = bind(check_size, _1, 6);
string s = "hello";		
bool b1 = check6(s);	//check(6)会调用check_size(s, 6);
auto wc = find_if(words.begin(), words.end(), bind(check_size, _1, sz));
//bind调用生成一个可调用对象，check_size第一个参数为words中的元素，第二个参数绑定到sz的值
```

#### 使用placeholders名字

名字`_n` 定义在 `placeholders` 命名空间中，而此命名空间定义在`std` 命名空间中，两个命名空间都要写上。placeholders命名空间也定义在头文件 `functional` 中。

```c++
using namespace std;
using namespace std::placeholders;
```

#### bind的参数

```c++
auto g = bind(f, a, b, _2, c, _1);
```

调用`g(X, Y) ` 实际会调用 `f(a, b, Y, c, X)` 。

```c++
//用bind重排参数序列
sort(words.begin(), words.end(), isShorter);
sort(words.begin(), words.end(), bind(isShorter, _2, _1));
```

#### 绑定引用参数

bind只能拷贝参数，如果想要传递给bind一个对象而又不拷贝它，使用标准库 `ref` 函数。

```c++
ostream &print(ostream &os, const string &s, char c) {
    return os << s << c;
}
for_each(words.begin(), words.end(), bind(print, ref(os), _1, ' '));
```

函数 `ref` 返回一个对象，包含给定的引用，此对象是可拷贝的。`cref` 函数生成一个保存const引用的类。两个函数定义在头文件 `functional` 中。

### 7. 插入迭代器

通过一个插入迭代器进行赋值时，该迭代器调用容器操作来向给定容器的指定位置插入一个元素。

`back_inserter` 创建一个使用push_back的迭代器。

`front_inserter` 创建一个使用push_front的迭代器。

`inserter` 创建一个使用insert的迭代器。此函数接受第二个参数，这个参数必须是一个指向给定容器的迭代器。元素将被插入到给定迭代器所表示的元素之前。（插入过程中，迭代器指向的元素不会变）

```c++
//it是由inserter生成的迭代器
*it = val;	//等价于以下两行代码
it = c.insert(it, val);	//插入后it指向新插入的元素
++ it;	//递增it使它指向原来的元素
```

```c++
list<int> lst = {1, 2, 3, 4};
list<int> lst2, lst3;		//空list
copy(lst.begin(), lst.end(), front_inserter(lst2));	//lst2为4 3 2 1
copy(lst.begin(), lst.end(), inserter(lst3, lst3.begin()));	//lst3为1 2 3 4
```



















































 

