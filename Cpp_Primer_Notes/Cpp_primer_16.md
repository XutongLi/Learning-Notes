## 模板与泛型编程

### 1. 泛型编程

泛型编程是独立于任何特定类型来编写代码的，在编译时可以获知类型。

### 2. 函数模板

一个函数模板就是一个公式，可用来生成针对特定类型的函数版本。

模板定义以 `template` 关键字开始，后跟一个模板参数列表（不能为空）。模板参数表示在类或函数定义中用到的类型或值。当使用模板时，指定模板实参，将其绑定到模板参数上。

模板中的函数参数是const的引用，这样保证了函数可以用于不能拷贝的类型，且能使函数运行的更快。

```c++
template <typename T>
int compare(const T &v1, const T &v2) {
    if (v1 < v2)	return -1;
    if (v2 < v1)	return 1;
    return 0;
}
```

T的实际类型在编译时根据compare的使用情况来确定。

编译器用推断出的模板参数来 **实例化** 一个特定版本的函数。这些编译器生成的版本通常被称为模板的 **实例** 。

如T被替换成int的实例为：

```c++
int compare(const int &v1, const int &v2) {
    if (v1 < v2)	return -1;
    if (v2 < v1)	return 1;
    return 0;
}
```

模板的头文件通常既包括声明也包括定义。

#### 模板类型参数

**类型参数** 可用来指定返回类型或函数的参数类型，以及在函数体内用于变量声明或类型转换。

类型参数前使用关键字 `typename` 或 `class` （两者等价，一般用typename）。

```c++
template <typename T>
T foo(T *p) {
    T *tmp = p;
    ...
    return tmp;
}
```

#### 非类型模板参数

可通过一个 **特定的类型名** 而非关键字class或typename来指定 **非类型参数** 。非类型模板参数的模板实参必须是 **常量表达式** 。

```c++
template<unsigned N, unsigned M>
int compare(const char (&p1)[N], const char (&p2)[M]) {
    return strcmp(p1, p2);
}
compare("hi", "mom");
//编译器实例化出的版本：
int compare(const char (&p1)[3], const char(&p2)[4])
```

### 3. 类模板

与函数模板不同的是，编译器不能为类模板推断模板参数类型，必须在模板名后的尖括号（ **显式模板实参** ）中提供额外信息，用来代替模板参数的模板时参列表。

```c++
template <typename T> class Blob {
public:
    Blob();
    void push_back(const T &t) { data->push_back(t); }
    ...
private:
    shared_ptr<vector<T>> data;
    ...
};
```

当编译器从Blob模板实例化出一个类时，会重写Blob模板，将模板参数T的每个实例替换为给定的模板实参。一个类模板的每个实例都形成一个独立的类。

#### 类模板的成员函数

与其他任何类相同，既可以在类模板内部，也可以在类模板外部为其定义成员函数。

```c++
//普通成员函数
template <typename T>
void Blob<T>::check(size_type i, const string &msg) const {
    if (i >= data->size())
        throw out_of_range(msg);
}
//构造函数
tempalte <typename T>
    Blob<T>::Blob() : data(make_shared<vector<T>>()) { }
```

成员函数只有在被用到时才进行实例化。

**注意** ：类模板的作用域内，可以直接使用模板名而不必指定模板实参：

```c++
template <typename T> 
class BlobPtr {
private:
    BlobPtr& operator++();	//前置运算符，等价于 BlobPtr<T>& operator++();
};
//后置运算符的类外定义
template <typename T>
BlobPtr<T> BlobPtr<T>::operator++(int) {
    BlobPtr ret = *this;	//等价于 BlobPtr<T> ret = *this
    ++ *this;
    return ret;
}
```

#### 类模板和友元

如果一个类模板包含一个非模板友元，则友元被授权可以访问所有模板实例。

如果友元自身是模板，类可以授权给所有友元模板实例，也可以只授权给特定实例。

```c++
template <typename T> class Pal;	//前置声明，将模板的一个特定实例声明为友元时用到
class C {	//C是一个非模板类
    friend class Pal<C>;	//用C实例化的Pal是C的一个友元
    //Pal2的所有实例都是C的友元，这种情况无需前置声明
    template <typename T> friend class Pal2;	
};
template <typename T> class C2 {	//C2是一个模板类
    //C2的每个实例将相同实例化的Pal声明为友元
    friend class Pal<T>;	
    //Pal2的所有实例都是C2每个实例的友元，不需要前置声明
    template <typename X> friend class Pal2;
    //Pal3是一个非模板类，Pal3是C2所有实例的友元
    friend class Pal3;
};
```

为了让所有实例成为友元，友元声明中必须使用与类模板本身不同的模板参数。

可令模板自己的类型参数成为友元：

```c++
template <typename T> 
class Bar {
friend T;
};
```

#### 类模板的static成员

```c++
template <typename T>
class Foo {
public:
    static size_t count() { return ctr; }
private:
    static size_t ctr;
};
Foo<string> fs;	//实例化static成员Foo<string>::ctr和Foo<string>::count
Foo<int> a, b, c;	//三个对象共享Foo<int>::ctr和Foo<int>::count
```

每个类模板实例都有自己的static成员实例。

static **数据成员** 也要定义为模板：

```c++
template <typename T>
size_t Foo<T>::ctr = 0;
```

与其他成员函数相同，static **成员函数** 只有在使用时才会实例化。

### 4. 模板参数

#### 模板声明

```c++
template <typename T> class Blob;
template <typename T> T calc(const &T, const &T);
template <typename U> U calc(const &U, const &U);	//两个calc等价
```

一个特定文件所需要的所有模板的声明通常一起放置在文件开始位置，出现于任何使用这些模板的代码之前。

#### 使用类的类型成员

类的类型成员是定义在类中的类。

默认情况下，C++通过作用域运算符访问的名字不是类型。因此，如果希望使用一个模板类型参数的类型成员，必须显式告诉编译器该名字是一个类型。通过使用关键字 `typename` 来实现这一点：

```c++
template <typename T>
typename T::value_type top(const T &c) {
    if (!c.empty())
        return c.bacK();
    else
        return typename T::value_type();	//返回值初始化的元素
}
```

#### 模板默认实参

##### 函数模板的默认实参

```c++
//compare有一个默认模板实参less<T>和一个默认函数实参F()
template <typename T, typename F = less<T>>
int compare(const T &v1, const T &v2, F f = F()) {
	if (f(v1, v2))	return -1;
	if (f(v2, v1))	return 1;
	return 0;
}
bool i = compare(1, 2);	//使用默认的less
Sales_data item1("111"), item2("222");
bool j = compare(item1, item2, compareIsbn);	
//第三个实参是一个可调用对象，该可调用对象的返回类型必须能转换为bool值。
```

类型参数F表示可调用对象的类型，默认模板实参为less\<T>。f为函数形参，默认值为F()。

##### 类模板的默认实参

如果一个类模板为其所有模板参数都提供了默认实参，要是用这些默认实参，就必须在模板名之后跟一个空尖括号对：

```c++
template <class T = int> 
class Numbers {
public:
	Numbers(T v = 0) :val(v) {}
private:
	T val;
};
Numbers<double> preci;
Numbers<> preci2;	//空<>表示希望使用默认类型
```

### 5. 成员模板

一个类（无论是普通类还是类模板）可以包含本身是模板的成员函数。这种成员称为 **成员模板** 。

成员模板不能是虚函数。

#### 普通类的成员模板

```c++
class DebugDelete {
public:
    DebugDelete(ostream &s = cerr) : os(s) {}
    template <typename T> void operator() (T *p) const {
        os << "deleting ptr" << endl;
        delete p;
    }
private:
    ostream &os;
};

double * p = new double;
DebugDelete d;
d(p);	//调用DebugDelete::operator()(double*)释放p
int *a = new int;
DebugDelete()(a);	//在一个临时DebugDelete对象上调用operator()(int*)
```

可将DebugDelete作为 `unique_ptr` 的删除器，在尖括号内给出删除器类型，并此类型对象给unique_ptr的构造函数：

```c++
unique_ptr<int, DebugDelete> p(new int, DebugDelete());
//销毁p指向的对象
//实例化DebugDelete::operator()<int>(int *)
```

#### 类模板的成员模板

对与类模板定义成员模板，类和成员各自有自己的独立的模板参数。

在类模板外定义成员模板时，必须同时为类模板和成员模板提供模板参数列表：

```c++
template <typename T> class Blob {
public:
	template <typename It> Blob(It b, Ib e);  
};
//定义：
template <typename T>	//类的类型参数
template <typename It>	//构造函数的类型参数
Blob<T>::Blob(It b, It e) : data(make_shared<vector<T>>(b, e)) {}
```

为了实例化类模板的成员模板，必须同时提供类和函数模板的实参：

```c++
int ia[] = {0, 1, 2, 3};
vector<long> vi = {0, 1, 2};
list<const char*> w = {"now", "is"};
//实例化Blob<int>类及其接受两个int*参数的构造函数，a1实例化为：Blob<int>(int*, int*)
Blob<int> a1(begin(ia), end(ia));	
//实例化Blob<int>类及其接受两个vector<long>::iterator参数的构造函数
Blob<int> a2(vi.begin(), vi.end());
//实例化Blob<string>类及其接受两个list<const char*>::iterator参数的构造函数
Blob<string> a3(w.begin(), w.end());
```

### 6. 控制实例化

相同实例可能出现在多个对象文件中，在多个文件中实例化相同模板的额外开销可能非常严重。可以通过 **显式实例化** 来避免这种开销。

```c++
extern template declaration;	//实例化声明
template declaration;			//实例化定义
```

将一个实例化声明为 `extern` 就表示承诺在程序其他位置有该实例化的一个非 extern 声明（定义）。对于一个给定的实例化版本，可能有多个extern声明，但只有一个定义。

extern声明必须出现在任何使用此实例化版本的代码之前。

```c++
//Application.cpp
//这些模板类型必须在程序其他位置进行实例化
extern template class Blob<string>;
extern template int compare(const int&, const int&);	//两个声明
Blob<string> s1, s2;	//实例化会出现在别的位置
Blob<int> a1 = {0, 1, 2};	//Blob<int>及其接受Initializer_list的构造函数在本文件中实现
int i = compare(a1[0], a1[1]);	//实例化出现在其他位置
```

```c++
//tempalteBuild.cpp
//实例化文件必须为每个在其他文件中声明为extern的类型和函数提供一个定义
template class Blob<string>;
template int compare(const int &, const in&);
```

编译时，将Application.o和tempalteBuild.o链接在一起。

在一个类模板的显式实例化定义中，所用类型必须能用于模板的所有成员函数。

### 7. 模板实参推断

对于 **函数模板** ，编译器利用调用中的函数实参来确定其模板参数。从函数实参来确定模板实参的过程为 **模板实参推断** 。

#### 类型转换与模板类型实参

将实参传递给 **模板类型** 的函数形参时，能够自动应用的类型转换只有：

- const转换，可以将一个非const对象的引用（或指针）传递给一个const的应用（或指针）形参。
- 数组或函数指针转换：如果函数形参 **不是引用类型** ，则可以对数组或函数类型的实参应用正常的指针转换。一个数组实参可以转换为一个指向其首元素的指针；一个函数实参可以转换为一个该函数类型的指针。

其他的类型转换，如算术转换、派生类向基类的转换以及用户定义的转换都不能应用于函数模板。

**注意** ：如果函数参数类型不是模板参数，则对实参进行正常的类型转换。

```c++
//ostream是非模板参数类型，可以正常类型转换；T是模板参数类型
template <typename T> ostream &print(ostream &os, const T &obj) {
    return os << obj;
}
```

#### 函数模板显式实参

某些情况下，编译器无法推断出模板实参的类型，如：

```c++
template <typename T1, typename T2, typename T3>
T1 sum(T2, T3);
```

编译器无法推断T1，它未出现在函数参数列表中。因此每次调用sum都必须为T1提供一个 **显式模板实参** 。显式模板实参在尖括号中给出，位于函数名之后，实参列表之前。

```c++
//T1是显式指定的，T2和T3是从函数实参类型推断而来的
auto val3 = sum<long long>(i, lng);	//long long sum(int, long)
```

显式模板实参按由左至右的顺序与模板参数匹配，尾部参数的现实模板实参可以忽略。

##### 显式实参的正常类型转换

```c++
compare<long>(lng, 1024);	//实例化compare(long, long)
compare<int>(lng, 1024);	//实例化compare(int, int)
```

#### 尾指返回类型与类型转换

对于不知道返回结果类型的情况，也可以使用尾置返回类型：

```c++
template <typename It>
auto fcn(It beg, It end) -> decltype(*beg) {
	//...
	return *beg;
}
```

通知编译器fcn的返回类型与解引用beg参数的结果类型形同，返回值是一个左值。











