## 关联容器

### 1. 关联容器与顺序容器区别

关联容器中的元素是按关键字来保存和访问的；顺序容器中的元素按它们在容器中的位置来顺序保存和访问。

### 2. 关联容器类型

| 类型               | 描述                             | 是否有序 | 头文件        |
| ------------------ | -------------------------------- | -------- | ------------- |
| map                | 关联数组，保存关键字-值对        | 有       | map           |
| set                | 关键字即值，即只保存关键字的容器 | 有       | set           |
| multimap           | 关键字可重复出现的map            | 有       | map           |
| multiset           | 关键字可重复出现的set            | 有       | set           |
| unordered_map      | 用hash函数组织的map              | 无       | unordered_map |
| unordered_set      | 用hash函数组织的set              | 无       | unordered_set |
| unordered_multimap | hash组织的map，关键字可重复出现  | 无       | unordered_map |
| unordered_multiset | hash阻止的set，关键字可重复出现  | 无       | unordered_set |



