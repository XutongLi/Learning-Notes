## 面向对象程序设计

### 1. 面向对象三个基本概念

**数据抽象** ：将类的接口与实现分离

**继承** ：可以定义相似的类型并对其相似关系建模

**动态绑定** ：可以在一定程度上忽略相似类型的区别，而以统一的方式使用它们的对象

### 2. 继承

**基类** 负责定义在层次关系中所有类共同拥有的成员，而每个 **派生类** 定义各自特有的成员。

派生类使用 **类派生列表** 指明它是从哪个（些）基类继承而来。

对于某些函数，基类希望它的派生类各自定义适合自身的版本，此时基类就将这些函数声明成 **虚函数** 。派生类必须在其内部对所有重新定义的虚函数进行声明（在该函数形参列表后增加一个 `override` 关键字）。

```c++
class Quote {
public:
    string isbn() const;
    virtual double net_price(size_t n) const;
};
class Bulk_quote : public Quote {	//基类前可以有访问说明符
public:
    double net_price(size_t n) const override;
};
```

### 3. 动态绑定

使用基类的 **引用或指针** 调用一个虚函数时将发生动态绑定。使用动态绑定，可用同一段代码分别处理 Quote 和 Bulk_quote 的对象。

```c++
//根据传入的item形参的对象类型调用Quote::net_price或Bulk_quote::net_price
double print_total(ostream &os, const Quote &item, size_t n) {
    double ret = item.net_price(n);
    os << ret << endl;
    return ret;
}
print_total(cout, basic, 20);	//调用Quote的net_price
print_total(cout, bulk, 20);	//调用Bulk_quote的net_price
```

### 4. 派生类

派生类通过 **类派生列表** 指明它是从哪个类继承而来。派生类可有多于一个基类。

一个派生类对象包含多个 **组成部分** ：一个含有派生类自己定义的（非静态）成员的子对象；以及一个与该派生类继承的基类对应的子对象，如果有多个基类，这样的子对象也有多个。

可以把派生类的对象当成基类对象来使用，而且也可以将基类的指针或引用绑定到派生类对象的基类部分上。

```c++
Quote item;			//基类对象
Bulk_quote bulk;	//派生类对象
Quote *p = &item;	//p指向Quote对象
p = &bulk;			//p指向bulk的quote部分
Quote &r = bulk;	//r绑定到bulk的quote部分
```

#### 派生类构造函数

派生类必须使用基类的构造函数来初始化基类部分（从基类中继承而来的成员）。

```c++
Bulk_quote(const string &book, double p, size_t qty, double disc) :
	Quote(book, p), min_qty(qty), discount(disc) {}
```

#### 继承与静态成员

若基类定义了静态成员，则在整个继承体系中只存在该成员的唯一定义。静态成员遵循通用的访问控制规则。

```c++
class Base {
public:
    static void statmem();
};
class Derived : public Base {
    void f(const Derived&);
};
void Derived::f(const Derived &obj) {
    //下列四种使用方法等价
    Base::statmem();		//Base定义了statmem
    Derived::statmem();		//Derived继承了statmem
    obj.statmem();			//通过Derived对象访问
    statmem();				//通过this对象访问
}
```

#### 防止继承的发生

要定义不能被继承的类，在其后加 `final` 关键字。

```c++
class A final {};	//A不能作为基类
class B final : public c {};	//B不能作为基类
```

### 5. 类型转换与继承

可以将基类的指针或引用绑定到派生类的对象上。

智能指针也支持派生类向基类的类型转换，可以将一个派生类对象的指针存储在一个基类的智能指针内。

派生类向基类的自动类型转换只对指针或引用类型有效。但是继承体系中大多数类定义了拷贝控制成员，因此可以将一个派生类对象拷贝、移动或赋值给一个基类对象（只处理派生类对象的基类部分）。

不存在从基类向派生类的 **隐式** 类型转换，但是在确定安全性的前提下使用 `static_cast` 将基类转换为派生类。

### 6. 静态类型和动态类型

表达式的 **静态类型** 在编译时总是已知的，它是变量声明的类型或表达式生成的类型；**动态类型** 则是变量或表达式表示的内存中的对象的类型，动态类型直到运行时才可知。

如：

```c++
auto ret = item.net_price(n);
```

item的静态类型时Quote&，动态类型依赖于item绑定的实参。基类的指针或引用的静态类型可能与其动态类型不一致。

如果表达式既不是引用也不是指针，则它的动态类型永远与静态类型一致。

### 7. 虚函数

**虚函数** 是基类希望派生类直接继承而不要改变的函数。当使用 *指针或引用* 调用虚函数时，该调用将被动态绑定，根据引用或指针所绑定的对象类型不同，该调用可能执行基类的版本，也可能执行某个派生类的版本。

虚函数声明语句之前加 `virtual` 关键字，virtual只能出现在类内部的声明语句，不能用于类外部的函数定义。

构造函数之外的非静态函数都可以是虚函数。

对虚函数的调用可能在运行时才被解析。

```c++
double print_total(ostream &os, const Quote &item, size_t n) {}
Quote base("aaa", 50);
print_total(cin, base, 10);	//调用Quote::net_price
Bulk_quote derived("bbb", 50, 5, 0.1);
prnt_total(cout, derived, 10);	//调用Bulk_quote::net_price
```

#### 派生类中的虚函数

一个函数被声明成虚函数，则它在所有派生类当中都是虚函数。声明可不加 `virtual` 关键字；在形参列表（包括const和引用修饰符）后加 `override` 关键字。

若将某个函数指定为 `final` ，则该函数不能被覆盖。

```c++
struct B {
	virtual void f() const;
};
struct B1 : B {
    void f() const override;
};
```

#### 回避虚函数的机制

有些情况下，希望对虚函数的调用不要执行动态绑定，而是强迫其执行虚函数的某个特定版本，使用作用于运算符可以实现这一目的。

```c++
//强行调用基类中定义的函数版本而不管baseP的动态类型是什么
double a = baseP->Quote::net_price(42);
```

通常情况下，只有成员函数（或友元）中的代码才需要使用作用域运算符来回避虚函数的机制。

一般在派生类的虚函数调用其覆盖的基类的虚函数版本时，使用此方法。

### 8. 多态性

当且仅当通过 **指针或引用** 调用虚函数时，才会在运行时解析该调用，这种情况下才有多形态性（对象的动态类型与静态类型不同）

对非虚函数的调用在编译时绑定，通过对象进行的函数调用（虚函数或非虚函数）也在编译时绑定。

### 9. 纯虚函数与抽象基类

#### 纯虚函数

在函数体的位置书写 `=0` 可以将一个虚函数说明为纯虚函数。=0只能出现在类内部的虚函数声明语句处。纯虚函数无需定义，如要定义只能定义在类的外部。

```c++
class Disc_quote : public Quote {
public:
    ...
    double net_price(size_t) const = 0;
};
```

#### 抽象基类

含有（或者未经覆盖直接继承）纯虚函数的类时 **抽象基类** 。抽象基类负责定义接口，而后续的其他类可以覆盖该接口。不能直接创建一个抽象基类的对象。

### 10. 重构

**重构** 负责重新设计类的体系以便将操作和/或数据从一个类移动到另一个类当中。