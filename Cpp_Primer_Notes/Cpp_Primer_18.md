## 用于大型程序的工具

### 1. 命名空间

不同的库中会有命名冲突，将多个库名字放置在全局命名空间中将引发 **命名空间污染** 。

**命名空间** 防止名字冲突，它分割了全局命名空间，其中每个命名空间是一个作用域。不同命名空间可以有相同名字的成员。

#### 命名空间定义

`namespace` 关键字加命名空间的名字。能出现在全局作用域中的声明就能置于命名空间内，主要包括：类、变量（及其初始化操作）、函数（及其定义）、模板和其他命名空间。

```c++
namespace cpp_primer {
    class Sales_data {};
    Sales_data operator+(const Sales_data&, const Sales_data&);
    class Query {};
}	//不需要分号结束
//位于该命名空间之外的代码必须明确指出所用的名字属于哪个命名空间
cpp_primer::Query q = cpp_primer:;Query("hello");
```

#### 命名空间可以是不连续的

```c++
namespace nsp {}
```

可能是定义了一个名为nsp的新命名空间，也可能是为已经存在的命名空间添加新成员。

这一特性可以将几个独立的接口和实现文件组成一个命名空间；可以将声明放在头文件中，定义放在源文件中。

相同命名空间中的不同类应该同不同的文件分别表示。

```c++
//Sales_data.h
#include <string>
namespace cpp_primer {
    class Sales_data {};
    Sales_data operator+(const Sales_data&, const Sales_data&); 	//.h文件中写入声明
}
//Sales_data.cpp
namespace cpp_primer {
    //Sales_data类的成员和重载运算符的定义			//.cpp文件中写入定义
}
//user.cpp
#include "Sales_data.h"
int main() {
    using cpp_primer::Sales_data;
    Sales_data a;	//或者 cpp_primer::Sales_data a;
    ...
    return 0;
}  
```

头文件一般不放在命名空间内部。如果这么做了，隐含的意思是把头文件所有的名字定义成该命名空间的成员。

#### 模板特例化

模板特例化必须定义在原始模板所属的命名空间中。和其他命名空间类似，只要在命名空间中声明了特例化，就能在命名空间外部定义它了。

```c++
namespace std {
    template <> struct has<Sales_data>;
}
//在std中添加了模板特例化的声明之后，就可以在命名空间std的外部定义它了
template <> struct std::hash<Sales_data> {
	...S  
};
```

#### 全局命名空间

全局作用域中定义的名字（即在所有类、函数及命名空间之外定义的名字）是定义在全局命名空间。全局作用域中定义的名字被隐式地添加到全局命名空间中。

`::member_name` 表示全局命名空间中的一个成员。

#### 嵌套命名空间

```c++
namespace cpp_primer {
    namespace QueryLib {
        class Query {};
    }
}
cpp_primer::QueryLib::Query a("hello");
```

#### 内联命名空间

与普通嵌套命名空间不同，**内联命名空间** 中的名字可以被外层命名空间直接使用。

定义内联命名空间的方式是在 `namespace` 前加 `inline` 。

```c++
namespace cpp_primer {
    intline namespace Fifth {
        class Query {};
    }
}
cpp_primer::Query a("hello");
```

关键字inline必须出现在命名空间第一次定义的地方，后续再打开命名空间的时候可以写inline，也可以不写。

#### 命名空间别名

```c++
namespace primer = cpp_primer;
namespace Qlib = cpp_primer::QueryLib;
Qlib::Query q;
```

声明别名必须在命名空间定义之后。

一个命名空间可以声明多个别名，它们等价。

#### using 声明

一条 **using 声明** 语句一次只引入命名空间的一个成员。

using 声明语句可以出现在全局作用域、局部作用域、命名空间作用域以及类作用域中。在类的作用域中，这样的声明语句只能指向基类成员。

当为函数写using声明的时候，该函数的所有版本都被引入到当前作用域中。

```c++
#include "Sales_data.h"
int main() {
    using cpp_primer::Sales_data;
    Sales_data a;	//或者 cpp_primer::Sales_data a;
    ...
    return 0;
}  
```

#### using 指示

`using 指示` 使某个特定的命名空间中所有的名字可见。using指示一般被看做是出现在最近的外层作用域中。

using 指示可以出现在 全局作用域、局部作用域和命名空间作用域中，但不能出现在类的作用域中。

```c++
namespace A {
    int i, j;
}
void f() {
    using namespace A;	//把A的名字注入到全局作用域中（即本函数作用域最近的外层作用域）
    					//（若全局作用域中已有i，此函数作用域中调用i会有二义性）
    cout << i * j << endl;	//使用命名空间A中的i和j
}
```

尽量避免使用 using 指示。但在命名空间本身的实现文件中可以使用 using 指示。

using 声明引入与已有函数形参列表完全相同的同名函数会引发错误，而 using 指示不会。

### 2. 多重继承

**多重继承** 是指从多个直接基类中产生派生类的能力。多重继承继承了所有父类的属性。

```c++
class Bear : public ZooAnimal {};
class Panda : public Bear, public Endangered {};	//派生列表中包含多个基类
```

多重继承关系中，派生类的对象含有每个基类的子对象。

派生类构造函数初始化所有基类：

```c++
//构造一个派生类的对象将同时构造并初始化它的所有基类子对象
Panda::Panda(string name, bool onExhibit) : Bear(name, onExhibit, "Panda"),
			Endangered(Endangered::critical) {}
```

### 3. 虚继承

默认情况下，派生类中含有继承链上每个类对应的子部分。如果某个类在派生过程中出现了多次，则派生类中将包括该类的多个子对象。

通过 **虚继承** 解决该问题。虚继承的目的是令某个类做出声明，承诺愿意共享它的基类。其中，共享的基类子对象称为 **虚基类** 。在这种机制下，不论虚基类在继承体系中出现了多少次，在派生类中都只包含一个共享的虚基类子对象。

#### 使用虚基类

指定虚基类的方式是在派生列表中添加关键字 `virtual` ：

```c++
//将ZooAnimal定义为Raccoon和Bear的虚基类
class Raccoon : public virtual ZooAnimal {};
class Bear : virtual public ZooAnimal {};	//virtual和public的顺序任意
```

在后续的派生类中共享虚基类的同一份实例：

```c++
class Panda : public Bear, public Raccoon, public Endangered {};
//因为Raccoon和Bear继承ZooAnimal的方式都是虚继承，所以在Panda中只有一个ZooAnimal基类部分
```

#### 构造函数与虚继承

虚派生中，虚基类由最低层的派生类初始化。

```c++
Panda::Panda(string name, bool, onExhibit):
	ZooAnimal(name, onExhibit, "Panda"),	
	Bear(name, onExhibit), Raccoon(name, onExhibit),
	Endangered(Endangered::cirtical), sleeping_flag(false) {}
//Panda位于派生的最低层并负责初始化共享的ZooAnimal基类部分
```

虚基类总是先于非虚基类构造，与它们在继承体系中的次序和位置无关。若有多个虚基类，虚基类的子对象按照它们在派生列表中出现的顺序从左到右依次构造。

对象的销毁顺序与构造顺序相反。



