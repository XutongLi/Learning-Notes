## 类

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

**区分** ：`const string getName()` 表示返回值是const类型。

### 4. 构造函数

构造函数的任务是初始化类对象的数据成员，类的对象被创建，就会执行构造函数。

构造函数的名字和类名相同，构造函数没有返回类型。类可以包含多个构造函数，和其他重载函数差不多，不同的构造函数之间必须在参数数量或类型上有区别。

构造函数不能被声明为const的。

### 5. 默认构造函数 

**默认构造函数** 无需任何实参。如果类没有显示地定义构造函数，编译器就会为我们隐式地定义一个默认构造函数。

如果一个构造函数为所有参数都提供了默认实参，则它实际上也定义了默认构造函数。

如果需要默认的行为，可以通过在参数列表后面写上`=default` 来要求编译器生成构造函数。其中`=default` 既可以和声明一起出现在类的内部（内联），也可以作为定义出现在类的外部（不是内联）。

```c++
struct Sales_data {
    Sales_data() = default;
}
```

默认构造函数的使用：

```c++
Sales_data obj();		//错误，声明了一个返回值类型为Sales_data的函数
Sales_data obj2;		//正确，调用默认构造函数生成了一个Sales_data对象
```

### 6. 构造函数初始值列表

```c++
struct Sales_data {
    Sales_data(cosnt string &s, usingned n, double p) : 
		bookNo(s), units_sold(n), revenue(p * n) {}
}
```

- 花括号定义了一个空的函数体。冒号以及冒号和花括号之间的代码称为**构造函数初始值列表** 。

- 当某个数据成员被构造函数的初始值列表忽略时，它将以与合成默认构造函数相同的方式隐式初始化。

- 如果成员是const、引用，或者属于某种未提供默认构造函数的类类型，我们必须通过构造函数初始值列表为这些成员提供初值。
- 成员的初始化顺序与他们在类定义中的出现顺序一致。构造函数初始值列表中初始值的前后位置关系不会影响实际的初始化顺序。

### 7. 委托构造函数

```c++
class Sales_data {
    public:
    	Sales_data(string s, unsigned cnt, double price) : 
    		bookNo(s), uints(cnt), rev(cnt * price) {}; 
    	//下列构造函数全部委托给了第一个构造函数
    	Sales_data() : Sales_data("", 0, 0) {};
   		Sales_data(string s) : Sales_data(s, 0, 0) {};
    	Sales_data(istream &is) : Sales_data() { read(is, *this) };
}
```

### 8. struct和class区别

使用class和struct定义类的**唯一** 区别就是默认的访问权限。

如果使用`struct` 关键字，则定义在第一个访问说明符之前的成员是public的；如果使用`class` 关键字，则这些成员是private的。

### 9. 友元

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

### 10. 令类成员作为内联函数

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

### 11. 可变数据成员

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

### 12. 返回*this的成员函数

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

### 13. 类的声明

```c++
class Screen;		//Screen类的声明
```

仅声明类而暂时不定义它，称为 **前向声明** ，Screen类在声明后定义前是一个 **不完全类型** 。

不完全类型的使用场景：可以定义指向这种类型的指针或引用，也可以声明（但不能定义）以不完全类型作为参数或者返回类型的函数。

因为只有当类全部完成后类才算被定义，所以一个类的成员类型不能是该类自己。但是类允许包含指向它自身类型的指针或引用。

### 14. 隐式的类类型转换

- 能通过 **一个实参** 调用的构造函数（**转换构造函数**）定义一条从构造函数的参数类型向类类型隐式转换的规则。

    ```c++
    class Sales_data {
        public:
            Sales_data() = default;
            Sales_data(const string &s) : bookNo(s) {};
            Sales_data &combine(const Sales_data&);
            ...
    }
    Sales_data item;
    string null_book = "9-9999-9999";
    item.combine(null_book);	//string被隐式转换为Sales_data临时对象
    ```

- 编译器只会自动地执行一步类型转换：

    ```c++
    item.combine("9-9999-9999");	//错误，发生了两部转换：字面值->string->Sales_data
    item.combine(Sales_data("9-9999-9999"));	//正确，隐式转换为string，显式转换为Sales_data
    item.combine(string("9-9999-9999"));	//正确，显式转换为string，隐式转换为Sales_data
    ```

- 将构造函数声明为 `explicit` 可以阻止构造函数定义的隐式转换：

    ```c++
    class Sales_data {
        public:
            Sales_data() = default;
            explicit Sales_data(const string &s) : bookNo(s) {};
            Sales_data &combine(const Sales_data&);
            ...
    };
    Sales_data Sales_data:(const string &s) {...}
    ```
    
    关键字`explicit` 只能对一个实参的构造函数有效。只能在类内声明构造函数时使用 `explicit` 关键字，在类外部定义时不应重复。
    
- 显式地转换类型（explicit只能阻止隐式转换）

    ```c++
    item.combine(Sales_data(null_book));
    item.combine(static_cast<Sales_data>(null_book));
    ```

### 15. 类的静态成员

#### 概念

类的静态成员与类本身直接相关，声明前加 `static` 关键字。静态成员可以是public或者private的。静态数据成员的类型可以是常量、引用、指针、类类型等。

```c++
class Account {
    public :
    	static double rate() { return amount += amount * inter ;}
    	static void rate(double);
    	...
    private:
    	double amount;
    	static double inter;
    	static double initRate();
};
```

类的静态成员存在于任何对象之外，对象中不包含任何与静态数据成员有关的数据。

静态成员函数不与任何对象绑定在一起，它不包含this指针。

#### 使用类的静态成员

```c++
double r = Account::rate();		//使用作用域运算符直接访问
Account ac1;
Account *p = &ac1;
r = ac1.rate();
r = p->rate();					//使用类的对象、引用或指针访问
```

#### 定义静态成员

**静态成员函数**

静态成员函数可以在类的内部和外部定义。

当在类的外部定义静态成员时，不能重复 `static` 关键字，该关键字只出现在类内部的声明语句。

**静态数据成员**

静态数据成员不能在类的内部初始化，必须在类的外部定义和初始化每个静态数据成员。

和全局变量类似，静态数据成员定义在任何函数之外，一旦被定义，就将一直存在于程序的整个声明周期中。

```c++
double Account::inter = initRate();
```

#### 静态成员的类内初始化

通常情况下，类的静态数据成员不应该在类的内部初始化，但是可以为静态成员提供const整数类型的类内初始值，不过要求静态成员必须是字面值常量类型的constexpr。

```c++
class Account {
    private:
    	static constexpr int period = 30;	//period是常量表达式
    	double daily_tbl[period];
};
constexpr int Account::period;
```

即使一个常量静态数据成员在类的内部被初始化了，也应该在类的外部定义一下该成员。

#### 静态成员应用的特殊场景（普通成员不可用）

- 静态数据成员的类型可以是它所属的类型（静态成员可以是不完全类型）
- 可以使用静态成员作为默认实参

### 16. 不完全类型

已经声明了但尚未定义的类型。

不完全类型不能用于定义变量或者类的成员，但是定义指针和引用是合法的。











