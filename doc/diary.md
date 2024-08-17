### 2024-08-17
从今天开始，在git里面记录下开发z语言的一起想法，毕竟时间一长，就不记得当时是怎么想的了

从7月20号第一次提交之后，陆陆续续进行了65次commit(话说，我们会不会65岁才退休，当然工作估计是到不了)，不断添加了一些特性和功能，现在能用来做一个简单的crud应用了，当然只是一个玩具。

那为啥要搞一个编程语言，重复的造轮子，还是一个很原始的轮子呢？

* 主要学习，搞一个语言，能更好使用其他语言
* 愿景：做出一个真正能实用的语言


开始今天的流水

目前语言实现基本的变量(integer, float, string, array, hash等)赋值申明，if，while，else，return等语句，支持package，import 方式引用文件，语言级内置基本函数(len,push,execute,typeof,http_server等)，builtin内置常用函数(is_string,is_int等)。

今天想着扩展支持for 语句，虽然用while也可以实现，比如:

```
let i = 0;
let len = 10;
while (i < 10) {
  var_dump(i)
  i = i + 1
}
```
但是这个语法，有一定的问题

* i这个变量作用域过大，本来只是循环使用，结果在外面也可以使用
* 如果忘记 i = i + 1，就死循环了

这个列子又发现一下问题，为啥不能用 i++， 其他都有呀！再不济用 i += 1 也行。。。

不用是因为语法不支持，但是其他语言都有，我们也要支持一下，先搞 i += 1，这个容易点，我们之前已经支持 i >= 1 格式

实现过程：
https://github.com/pantingwen/zlang/commit/5af6b9988d180463f2eba9b8be78afe289d9429c
区别点在于， += 操作，需要更新变量的值，方便后续使用

顺便优化一下词法解析器，复用 newTokenWithTwoChar，很多token，先 >=, <=, +=, -= 都可以用
https://github.com/pantingwen/zlang/commit/893cc0f015f87fbfc10273fa1a635dfc0a901edb

再顺便把 -=,*=,/= 都支持一下吧
https://github.com/pantingwen/zlang/commit/152a69af2d342639b2a0d9fcb5460899feebb74f

支持了++， --
https://github.com/pantingwen/zlang/commit/718735219b6f8b3fa4ede4d7b166bc923e89b31d