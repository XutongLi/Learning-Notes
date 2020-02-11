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















