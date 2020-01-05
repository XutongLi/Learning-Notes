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

`capture list` （捕获列表）是一个lambda所在函数中定义的局部变量的列表，如果为空表示此lambda不使用它所在函数中的任何局部变量。捕获列表只用于局部非static变量，lambda可以直接使用局部static变量和它所在函数之外声明的名字。

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
```



 

