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

