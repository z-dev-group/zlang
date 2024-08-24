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

处理if， while作用域的问题，同时 = 操作走中缀表达式，不要单独的表达式
https://github.com/pantingwen/zlang/commit/986939b341c30284a5a264fdccdf183a6cd2e40f
```
if (2 > 1) {
  let name = "seven"
}
var_dump(name) // 之前这里输出会是 seven 
```

### 2024-08-18

支持break关键字，在while语句中，遇到break跳出循环,后续for也将支持break

实现for语句
```
let count = 0

for (let i = 0; i < 10; i++) {
  count = i
}

var_dump(count) // expected 9

for (let i = 0; i < 10; i++) {
  count = i
  if (count > 5) {
    break
  }
}

var_dump(count) // expected 6
```
for 整体和while实现差不多，只是多个Initor，After，Initor是在程序循环前调用，After是每次body执行后调用，常用来修改条件

理论上可以共用，这里还是很开

接下来开始搞搞oop，目前程序hash支持函数的方式，就是对象方法

```
let dog = {
  "name": "siri"
  "say": fn(name) {
    var_dump(name + " is wangwang...")
  }
}

dog["say"](dog["name"])
```
添加语法糖，支持形如以下格式

```
class Dog {
  name: "siri",
  say: fn() {
    var_dump(this->name + "is wangwang...")
  }
}
dog = new Dog()
dog->say()
```

支持继承,继承最多两个吧，就像人继承自父母，实现接口支持一个
```
class ChinaDog extends Dog implement Intf
{
  from: "china"
}
```

```
interface Intf {
  say: fn() {}
  hello: fn() {}
}
```

### 2024-08-20
处理 oop 词法分析和抽象语法树定义


### 2024-08-22
这两天晚上陆陆续续做了 类class，接口interface，对象new的语法树的解析
支持-> 方式获取对象属性和方法

### 2024-08-24
初步实现对象操作
目前在new 一个对象的，会把对应class和父类的方法都拷贝一份,只支持let语句，就是函数也要 let 方式申明
面向对象要开发的东西还是挺多的，比如静态方法，私有，公开方法等等

* 支持直接函数方式
* 类访问 :: 直接访问类的数据 -> 是访问对象数据
* 函数可见性