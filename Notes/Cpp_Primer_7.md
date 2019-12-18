### 1. 类的基本思想

类的基本思想是**数据抽象** 和**封装** 。

**数据抽象** 是一种依赖于**接口** 和**实现** 分离的编程技术。

**封装** 实现了类的接口和实现的分离。

### 2. this指针

this指针是一个常量指针，成员函数可以通过this指针来访问调用它的那个对象。

```c++
string isbn() const {return this->bookNo;}	//isbn是一个类的成员函数
```

### 3. 常量成员函数

```c++
struct Sales_data{
    string bookNo;
    string isbn() const {return bookNo;}
};
```

常量成员函数在参数列表之后有const关键字。常量成员函数不能修改类数据成员、不能调用非常量函数。

它实际上是修改了**this** 指针的类型，将其修改为了指向常量对象的常量指针（原来是指向非常量对象的常量指针）。

### 4. 构造函数

构造函数的任务是初始化类对象的数据成员，类的对象被创建，就会执行构造函数。

构造函数的名字和类名相同，构造函数没有返回类型。类可以包含多个构造函数，和其他重载函数差不多，不同的构造函数之间必须在参数数量或类型上有区别。

构造函数不能被声明为const的。

**默认构造函数** 无需任何实参。如果类没有显示地定义构造函数，编译器就会为我们隐式地定义一个默认构造函数。

### 5. 默认构造函数 '=default' 的含义

如果需要默认的行为，可以通过在参数列表后面写上`=default` 来要求编译器生成构造函数。其中`=default` 既可以和声明一起出现在类的内部（内联），也可以作为定义出现在类的外部（不是内联）。

```c++
struct Sales_data {
    Sales_data() = default;
}
```

### 6. 构造函数初始值列表

```c++
struct Sales_data {
    Sales_data(cosnt string &s, usingned n, double p) : 
		bookNo(s), units_sold(n), revenue(p * n) {}
}
```

花括号定义了一个空的函数体。冒号以及冒号和花括号之间的代码称为**构造函数初始值列表** 。

当某个数据成员被构造函数的初始值列表忽略时，它将以与合成默认构造函数相同的方式隐式初始化。

### 7. struct和class区别

使用class和struct定义类的**唯一** 区别就是默认的访问权限。

如果使用`struct` 关键字，则定义在第一个访问说明符之前的成员是public的；如果使用`class` 关键字，则这些成员是private的。

### 8. 友元

类可以允许其他类或者函数访问它的非公有成员，方法是令其他类或者函数称为它的 **友元** 。

#### 友元函数

友元声明只能出现在类定义的内部，友元函数声明以 `friend` 关键字开始：

```c++
class Sales_data {
    friend Sales_data add(const Sales_data&, const Sales_data&);	//友元声明
    friend std::istream &read(std::istream&, Sales_data&);
    public:
    	...
    private:
    	...
}
Sales_data add(const Sales_data&, const Sales_data&);	//函数声明
std::istream &read(std::istream&, Sales_data&);
```

友元的声明仅仅指定了访问的权限，而非一个通常意义上的函数声明。需要在类内友元声明之外再专门对函数进行一次声明。

#### 友元类

```c++
class Screen {
    friend class Window_mgr;	//Window_mgr的成员可以访问Screen的私有成员
}
```

#### 令成员函数作为友元

```c++
class Screen {
    friend void Window_mgr::clear(index);	//Window_mgr类的clear函数可以访问Screen类的私有成员
}
```

声明和定义顺序：

- 首先定义Window_mgr类，其中声明clear成员函数，但不能定义它
- 定义Screen类，包括对与clear的友元声明
- 最后定义clear，此时它可以使用Screen的成员

### 9. 令类成员作为内联函数

```c++
class Screen {
    public:
    	...
    	char get() const {return contents[cursor]};		//隐式内联
    	inline char get(pos ht, pos wd) const;			//显式内联
    	Screen &move(pos r, pos c);						//能在之后被设为内联
	private:
    	...
}
char Screen::char get(pos r, pos c) const {...}
inline Screen &Screen::move(pos r, pos c) {...}	
```

定义在类内部的成员函数是自动inline的；可以在类的内部把inline作为声明的一部分显式地声明成员函数；也可以在类的外部用inline关键字修饰函数的定义。无须在声明和定义的地方同时说明inline。

### 10. 可变数据成员

将一个数据成员的声明中加入 `mutable` 关键字，这个数据成员的值可以被任何成员函数，包括const函数在内所改变。

```c++
class Screen {
    public:
    	void cnt() const;
    private:
    	mutable size_t access_ctr;
};
void Screen::cnt() const {
    ++ access_ctr;
}
```

### 11. 返回*this的成员函数

```c++
class Screen {
    public:
    Screen &set(char);
};
inline Screen Screen::set(char c) {
    contents[cursor] = c;
    return *this;	//将this指向的对象作为左值返回
}
```

返回引用的函数是左值的，意味着这些函数返回的是对象本身而非对象的副本。

```c++
myScreen.move(4, 0).set('#');	//把光标移到指定位置，然后设置该位置的字符值
myScreen.move(4, 0);
myScreen.set('#');				//等价
```

 注意：一个const成员函数如果以引用的形式返回*this，那么它的返回类型将是常量引用。

### 12. 类的声明

```c++
class Screen;		//Screen类的声明
```

仅声明类而暂时不定义它，称为 **前向声明** ，Screen类在声明后定义前是一个 **不完全类型** 。

不完全类型的使用场景：可以定义指向这种类型的指针或引用，也可以声明（但不能定义）以不完全类型作为参数或者返回类型的函数。

因为只有当类全部完成后类才算被定义，所以一个类的成员类型不能是该类自己。但是类允许包含指向它自身类型的指针或引用。





















