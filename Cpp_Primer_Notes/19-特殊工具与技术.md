## 特殊工具与技术

### 1. 运行时类型识别

#### 主要应用场景

使用基类指针或引用执行某个派生类操作，并且该操作不是虚函数。

#### 相关运算符

**运行时类型识别** （RTTI）的功能由两个运算符实现：

- `typeid` 运算符，用于返回表达式的类型
- `dynamic_cast` 运算符，用于将基类的指针或引用安全地转换成派生类的指针或引用

当将这两种运算符用于某种类型的指针或引用，并且该类型含有虚函数时，运算符将使用指针或引用所绑定对象的动态类型。

#### dynamic_cast运算符

执行从基类向派生类的带检查的强制类型转换。

```c++
dynamic_cast<type*>(e);		//e必须是一个有效指针，转换失败返回0
dynamic_cast<type&>(e);		//e必须是一个左值，转换失败抛出bad_cast异常
dynamic_cast<type&&>(e);	//e必须是一个右值，转换失败抛出bad_cast异常
```

**e** 的类型必须是：目标type的公有派生类、或目标type的公有基类、或目标type类。

##### 指针类型的dynamic_cast

```c++
//Derived是Base的派生类，bp指向对象的静态类型为Base
if (Derived *dp = dynamic_cast<Derived*>(bp)) {	//转换失败返回0
    //使用dp指向的Derived对象
}
else {
    //使用bp指向的Base对象
}
```

如果bp指向的是Derived对象，则上述类型转换初始化dp并令其指向bp所指的Derived对象。

在条件部分执行dynamic_cast操作可以确保类型转换和结果检查在同一条表达式中完成。

##### 引用类型的dynamic_cast

```c++
void f(const Base &b) {
    try {
        const Derived &d = dynamic_cast<const Derived&>(b);
    }
    catch (bad_cast) {
        //处理失败的情况
    }
}
```

对引用的类型转换失败时，程序抛出一个名为 `std::bad_cast` 的异常，该异常定义在 `typeinfo` 标准库头文件中。

#### typeid运算符

`typeid` 运算符用于获取对象的动态类型（当类型含有虚函数时）。

```c++
typeid(e);
```

e可以是任意表达式或者类型的名字，返回结果是一个常量对象的引用，该对象的类型是标准库类型 `type_info` 或 `type_info` 的公有派生类型。

通常情况下，使用typeid比较两条表达式的类型是否相同，或者比较一条表达式的类型是否与指定类型相同。

```c++
Derived *dp = new Derived;
Base *bp = dp;	//两个指针都指向Derived对象
//在运行时比较两个对象的类型
if (typeid(*bp) == typeid(*dp)) {
    //bp和dp指向同一类型的对象
}
if (typeid(*bp) == typeid(Derived)) {
    //bp实际指向Derived对象
}
```

typeid是否需要运行时检查决定了表达式是否会被求值。只有当类型含有虚函数时，编译器才会对表达式求值。反之，如果类型不含有虚函数，则typeid返回表达式的静态类型；编译器无需对表达式求值也能知道表达式的静态类型。

##### type_info类

`type_info` 类定义在 `typeinfo` 头文件中。

type_info类没有默认的构造函数，而且它的拷贝和移动构造函数以及赋值运算符都被定义成删除的。因此，无法定义或拷贝type_info类型的对象，也不能为type_info类型的对象赋值。创建type_info对象的唯一途径是使用`typeid` 运算符。

### 2. 类成员指针

**成员指针** 是指可以指向类的非静态成员的指针。

#### 数据成员指针

```c++
string Screen::*pdata;
pdata = &Screen::contents;	//将取地址运算符作用于Screen类的成员而非内存中的一个该类对象
auto pdata = &screen::contents;	//或者使用auto

Screen myScreen, *pScreen = &myScreen;
auto s = myScreen.*pdata;
s = pScreen->*pdata;	//或者
//首先解引用成员指针得到所需的成员，然后通过对象(.*)或指针(->*)获取成员
```

初始化一个成员指针时，该指针没有指向任何数据。只有当解引用成员指针时才提供对象的信息。

##### 返回数据成员指针的成员函数

```c++
class Screen {
public:
    //data是一个静态成员函数，返回一个成员指针
    static const std::string Screen::*data() {
        return &Screen::contents;
    }
};
//data返回一个指向Screen类的contents成员的指针
const string Screen::*pdata = Screen::data();
auto s = myScreen.*pdata;	//获取对象的成员
```

#### 成员函数指针

```c++
auto pmf = &Screen::get_cursor;
```

类似于其他函数指针，指向成员函数的指针也需要指定目标函数的返回类型和形参列表。如果成员函数时const成员或者引用成员，则我们必须将const限定符或引用限定符包含进来。

若成员存在重载的问题，则必须显式地声明函数类型，指明想要使用的是哪个函数：

```c++
char (Screen::*pmf2)(Screen::pos, Screen::pos) const;	
pmf2 = &Screen::get;	
```

##### 使用成员函数指针

```c++
//使用成员函数指针
Screen myScreen, *pScreen = &myScreen;
//通过pScreen所指的对象调用pmf所指的成员函数
char c1 = (pScreen->*pmf)();
//通过myScreen对象将实参0,0传给含有两个形参的get函数
char c2 = (myScreen.*pmf2)(0, 0);
```

### 3. 嵌套类

一个类可以定义在另一个类的内部，前者称为 **嵌套类** 。

嵌套类也使用访问限定符来控制外界对其成员的访问权限。

嵌套类必须声明在类的内部，但是可以定义在类的内部或是外部。嵌套类可以直接使用外层类的成员。

外层类的成员可以像使用其他任何类型成员一样使用嵌套类的名字。

```c++
class TextQuery {
public:
    class QueryResult;
};
//定义嵌套类
class TextQuery::QueryResult {	
	friend ostream& print(ostream&, const QueryResult);
public:
    QueryResult(string, int);
};
//定义嵌套类成员
TextQuery::QueryResult::QueryResult(string s, int i) : str(s), in(i) {}
```





