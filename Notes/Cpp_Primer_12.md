## 动态内存

### 1. 各对象生存期

**全局对象** 在程序启动时分配，在程序结束时销毁。

**局部自动对象** 在进入其定义所在的程序块时创建，在离开块时销毁。

**局部static对象** 在第一次使用前分配，在程序结束时销毁。

### 2. 动态对象释放

C++支持动态分配对象。动态对象在 **堆** 中分配空间。动态分配的对象的生存期与它们在哪里创建是无关的，只有当显式地被释放时，这些对象才会销毁。但是动态对象的正确释放是编程中极易出错的地方。

为了更安全地使用动态对象，标准库定义了两个 **智能指针** 类型来管理动态分配的对象。当一个对象应该被释放时，指向它的智能指针可以确保自动地释放它。

### 3. 使用动态内存原因

- 程序不知道自己需要多少对象（例如容器）
- 程序不知道所需要对象的准确类型
- 程序需要在多个对象间共享数据

### 4. 智能指针

**智能指针** 负责自动释放所指向的目标。

新标准库提供的两种智能指针的区别在于管理底层指针的方式：`shared_ptr` 允许多个指针指向同一个对象；`unique_ptr` 独占所指向的对象。`weak_ptr` 是伴随类，是一种弱引用，指向shared_ptr所指向的对象。

三种类型都定义在 `memory` 头文件中。

### 5. shared_ptr

#### make_shared

`make_shared` 标准库函数在动态内存中分配一个对象并初始化它，返回指向此对象的shared_ptr。

```c++
//指向一个值为42的int的shared_ptr
shared_ptr<int> p3 = make_shared<int>(42);
//指向一个值为"aaa"的string的shared_ptr
auto p4 = make_shared<string>(3, 'a');
```

make_shared用其参数来构造给定类型的对象；如果不传递任何参数，对象就会进行值初始化。

#### 引用计数

每个shared_ptr都有一个关联的计数器，称之为 **引用计数** 。

拷贝一个shared_ptr时，计数器会递增。如当用一个shared_ptr初始化另一个shared_ptr，或将它作为参数传递给下一个函数以及作为函数的返回值时。

当给shared_ptr赋予一个新值或是shared_ptr被销毁时，计数器会递减。

当一个shared_ptr的 *计数器变为0* ，它就会自动释放自己所管理的对象。

```c++
auto r = make_shared<int>(42);	//r指向的int只有一个引用者
r = q;	//给r赋值，令它指向另一个地址
		//递增q指向对象的引用计数
		//递减r原来指向对象的引用计数
		//r原来指向的对象已没有引用者，会自动释放
```

#### 自动销毁对象并释放内存

当指向一个对象的最后一个shared_ptr被销毁时，shared_ptr类会自动销毁此对象。它是通过 **析构函数** 完成销毁工作的。析构函数控制对象销毁时做什么操作。

shared_ptr还会自动释放相关联的内存。

#### get

智能指针类型定义了一个名为 `get` 的函数，它返回一个内置指针，指向智能指针管理的对象。

此函数为了这样一种情况设计：需要向不能使用智能指针的代码传递一个内置指针。

使用get返回的指针的代码不能delete此指针。

**注意** ：get用来将指针的访问权限传递给代码，只有在确定代码不会delete指针的情况下才可以使用get。特别是，永远不要用get初始化另一个智能指针或者为另一个智能指针赋值。

#### reset

可以用 `reset` 来将一个新的指针赋予一个shared_ptr

```c++
p.reset(new int(1024));	//p指向一个新对象
```

reset会更新引用计数，需要的话，会释放p指向的对象。

reset经常与 `unique` 一起使用，来控制多个shared_ptr共享的对象。在改变底层对象之前，检查自己是否是当前对象仅有的用户；若不是，改变之前要制作一份新的拷贝。

```c++
if (!p.unique())
    p.reset(new string(*p));	//不是唯一用户，分配新的拷贝
*p += 1;	//现在是唯一的用户，可以改变对象的值
```

### 6. 动态内存管理

`new` ：在动态内存中为对象分配空间并返回一个指向该对象的指针。

`delete` ：接受一个动态对象的指针，销毁该对象，并释放与之关联的内存。

#### new

```c++
int *pi = new int;	
//在自由空间构造一个int型对象，并返回指向该对象的指针
//pi指向一个动态分配的、未初始化的五名对象
```

默认情况下，动态分配的对象是默认初始化的，这意味着内置类型或者组合类型的对象的值将是 *未定义的* ，而类类型对象将用默认构造函数进行初始化。

初始化一个动态分配的对象，可以使用传统的构造方式（圆括号），也可以使用列表初始化（花括号）。

```c++
int *pi = new int(1024);
string *ps = new string(10, 'a');
vector<int> *pv = new vector<int>{0, 1, 2, 3 ,4};
```

#### 内存耗尽

自由空间被耗尽的情况是有可能发生的，如果 `new` 不能分配所要求的的内存空间，它会抛出一个类型为 `bad_alloc` 的异常。

可以使用 **定位new** 来阻止它抛出异常。定位new表达式允许我们向new传递额外的参数。

```c++
int *p1 = new int;	//如果分配失败，new抛出std::bad_alloc
int *p2 = new (nothrow) int;	//如果分配失败，new返回一个空指针
```

#### delete

通过 **delete表达式** 将动态内存归还给系统。delete表达式接受一个指针，指向我们想要释放的对象。它执行两个动作：销毁给定指针指向的对象；释放对应的内存。

```c++
delete p;
Foo *factory(T arg)	return new Foo(arg);
void use_factory(T arg) {
    Foo *p = factory(arg);
    delete p;	//若此处不释放内存，p离开了作用域之后，指向的内存就无法释放了
}
```

**注意** ：传递给delete的指针必须是指向动态分配的内存，或者是一个空指针。释放一块并非new分配的内存，或者将相同的指针值释放多次，其行为是未定义的。

在delete之后，指针就变成了 **空悬指针** ，即，指向一块曾经保存数据对象但现在已经无效的内存的指针。为初始化指针的所有缺点空悬指针也都有。可以在delete之后将nullptr赋予指针，这样就清楚地指出指针不指向任何对象。

#### 动态分配const对象

用new分配const对象是合法的。

```c++
const int *p = new const int(2014);	//分配并初始化一个const int
const string *ps = new const string;	//分配并默认初始化一个const的空string
```

一个动态分配的const对象必须进行初始化。对于一个定义了默认构造函数的类类型，其const动态对象可以隐式初始化，而其他类型的对象就必须显式初始化。

### 7. shared_ptr和new结合使用

```c++
shared_ptr<int> p1(new int(42));
shared_ptr<int> clone(int p)
    return shared_ptr<int>(new int(p));
```

接受指针参数的智能指针构造函数是 `explicit` 的，因此不能将一个内置指针隐式转换成一个智能指针，必须使用直接初始化形式来初始化一个职能指针。

#### 不能混用普通指针和智能指针

当将一个shared_ptr绑定到一个普通指针时，就将内存的管理责任交给了这个shared_ptr。

这样做了之后就不应该再使用内置指针来访问shared_ptr所指向的内存，因为不知道对象何时会被销毁。

```c++
int *x(new int(1024));
process(shared_ptr<int>(x));	//合法，但内存会被释放
int j = *x;	//未定义的，x是一个空悬指针
```

正确用法：

```c++
shared_ptr<int> p(new int(42));	//引用计数为1
process(p);	//引用计数为2
int i = *p;	//正确，引用计数值为1
```

### 8. shared_ptr管理非new分配内存的资源

定义删除器函数来完成对shared_ptr中保存的指针进行释放的操作。

```c++
void end_connection(connection *p)	{ disconnect(*p); }
void f(destination &d) {
    connection c = connect(&d);
    shared_ptr<connection> p(&c, end_connection);
    //使用连接
    //当f退出时（即使是由于异常退出），connection会被正确关闭
}
```

### 9. unique_ptr

某个时刻只能有一个unique_ptr指向一个给定对象。当unique_ptr被销毁时，它所指向的对象也被销毁。

与shared_ptr不同，没有类似make_shared的标准库函数返回一个unique_ptr。当定义一个unique_ptr时，需要将其绑定到一个new返回的指针上。unique_ptr不支持普通的拷贝或赋值操作。

```c++
unique_ptr<int> p2(new int(24));
```

可以通过调用`release` 或者 `reset` 将指针的所有权从一个（非const）unique_ptr转移给另一个unique_ptr：

```c++
unique_ptr<int> p2(p1.release());	//将所有权从p1转给p2，p1被release置为空
p2.reset(p3.release());	//将所有权从p3转给p2，reset释放了p2原来指向的内存
```

向unique_ptr传递删除器用法与shared_ptr有所不同

```c++
void f(destination &d) {
    connection c = connect(&d);
    unique_ptr<connection, decltype(end_connection)*> p(&c, end_connection);
    //使用连接
    //当f退出时（即使是由于异常退出），connection会被正确关闭
}
```

### 10. weak_ptr

`weak_ptr` 是一种不控制所指向对象生存期的智能指针，它指向由一个shared_ptr管理的对象。

将一个weak_ptr绑定到一个shared_ptr不会改变对象的引用计数。一旦最后一个指向对象的shared_ptr被销毁，对象就会被释放（即使有weak_ptr指向对象）。

当创建一个weak_ptr时，要用一个shared_ptr来初始化它：

```c++
auto p = make_shared<int>(24);
weak_ptr<int> wp(p);	//wp弱共享p，p的引用计数未改变
```

由于对象可能不存在，所以不可以使用weak_ptr直接访问对象，而必须调用lock。此函数调查weak_ptr所指向的对象是否存在，若存在，返回一个指向共享对象的shared_ptr。

```c++
if (shared_ptr<int> np = wp.lock()) {//若np不为空则条件成立
    //只有当lock返回true才进入if语句体，
    //if中，np与wp共享对象
}
```

























