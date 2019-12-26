## IO库

### 1. IO库

- `iostream` 处理控制台IO
- `fstream` 处理命名文件IO
- `stringstream` 完成内存string的IO 

### 2. IO类型间关系

标准库流特性可以无差别应用于普通流、文件流和string流，以及char或w_char版本。

类 `fstream` 和 `stringstream` 都继承自类 `iostream` 。

标准库使我们忽略不同类型的流之间的差异，这时通过 **继承机制** 实现的。

### 3. IO对象无拷贝或赋值

不能拷贝或对IO对象赋值，因此也不能将形参或者返回类型设置为流类型。

进行IO操作的函数通常以引用方式传递和返回流。读写一个IO对象会改变其状态，因此传递和返回的引用不能是const的。

### 4. rdstate、setstate和clear

流对象的 `rdstate` 成员返回一个 `iostate` 值，对应流的当前状态。

`setstate` 操作将给定条件位置位，表示发生了对应错误。

`clear` 不接受参数的版本清除（复位）所有错误标志位；接收参数的版本复位给定的条件状态位。

```c++
auto old_state = cin.rdstate();		
cin.clear();			
process_input(cin);				//使用cin
cin.setstate(old_state);		//将cin置为原有状态
//复位failbit和badbit，保持其他标志位不变
cin.clear(cin.rdstate() & ~cin.failbit & ~cin.badbit);
```

### 5. 刷新输出缓冲区

```c++
cout << "hi!" << endl;		//输出hi和一个换行，然后刷新缓冲区
cout << "hi!" << flush;		//输出hi，然后刷新缓冲区，不附加任何额外字符
cout << "hi!" << ends;		//输出hi和一个空字符，然后刷新缓冲区
```

### 6. unitbuf操纵符

```c++
cout << unitbuf;	//所有输出操作后都会立即刷新缓冲区
//任何输出都立即刷新，无缓冲
cout << nounitbuf;	//回到正常的缓冲方式
```

### 7. 如果程序崩溃，输出缓冲区不会被刷新

### 8. 关联输入和输出流

当一个输入流被关联到一个输出流时，任何试图从输入流读取数据的操作都会先刷新关联的输出流。标准库将cout和cin关联在一起。

`tie` 有两个重载版本：一个不带参数，如果本对象关联到一个输出流，则返回指向此流的指针，否则返回空指针；第二个版本接受一个指向 `ostream` 的指针，将自己关联到此ostream。

可以将一个istream对象关联到另一个ostream，也可以将一个ostream关联到另一个ostream。

```c++
cin.tie(&cout);
ostream *old_tie = cin.tie(nullptr);	//old_tie指向当前关联到cin的流，然后cin不再与其他流关联
cin.tie(&cerr);		//将cin与cerr关联，读取cin会刷新cerr而不是cout
cin.tie(old_tie);	//重建cin和cout之间的正常关联
```

每个流最多关联到一个流，但多个流可以同时关联到同一个ostream。

### 9. 阻止ofstream清空文件内容

```c++
//以语句会清空（即截断）f1内容
ofstream out("f1");
ofstream out2("f1", ofstream::out);
ofstream out3("f1", ofstream::out | ofstream::trunc);
//以下语句会保留f1内容（必须显式指定app模式）
ofstream app("f1", ofstream::app);	//隐含为输出模式
ofstream app2("f1", ofstream::out | ofstream::app);
```



