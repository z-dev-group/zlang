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

对于可见性这块，就支持私有，公开两种吧
_ 开头的变量为私有，否则公开, __ 是protectd 保护状态
通过对象或者类去访问属性方法，如果是_ 开头，直接返回不存在
对了，先实现对类方法的直接访问，就是 fn name() {} ----> let name = fn() {}， 支持了


基本实现面向对象的功能，还有一些问题需要解决
父类的属性，在对象中修改后，父类方法无法读取等

### 2024-09-12
又偷懒了，好久不更新
* 处理http_server调用对象方法时找不到父类方法的问题
* 支持rudev，更新z代码时，自动重新执行

### 2024-09-13
需要实现json_decode方法,但是在定义一个json字符串的时候，需要处理双引号的问题，我们需要支持``方式定义字符串
json_decode 实现比较简单，把参数作为代码运行一下
* 实现json_encode, json_decode

### 2024-09-17
中秋几天放假，基本就更新了一下z语言文档，地址是：[https://z-dev-group.github.io/zlang/](https://z-dev-group.github.io/zlang/)，语言功能没有更新，
要不做一下错误处理，go的实现是返回多个参数，用一个来保存错误信息，java使用try catch捕获异常
个人想法能不能统一用一个参数返回错误呢，通过with_error内置函数给一个普通变量携带上错误信息，通过is_with_error内置函数，判断变量是否错误信息，get_error_message获取变量错误信息
with_error(variable, errorMessageStr) // 设置错误信息
is_with_error(variable); // true or false
get_error_message(variable);  // 获取错误信息
* 错误处理,支持变量携带错误信息 ---> (float, integer, boolean, string, hash, array)

### 2024-09-18
z语言的hash类型，也就是object类型，使用的go的map，当时map是没有顺序的，json_encode的时候，顺序不固定，想要修复一下，顺序固定
* 使得hash类型json_encode顺序固定

### 2024-09-21
把hash统一叫做object，避免歧义

```
let person = {
  "name": "seven",
  "age": 12
}
typeof(person)
```

### 2024-09-22
支持syscall，使用syscall封装系统相关函数，比如getpid，getppid

### 2024-10-09
又是好一段时间没有开发了，今天开始记录点，做defer功能吧
国庆去北京玩了7天，也经历了24年的牛市到牛屎，牛时候感觉马上要财务自由了，z语言不用开发了，结果2天跌回来，还是好好提升自己吧
defer 是 go语言一个比较特殊的功能，可以在程序返回前执行，z语言模仿该功能，不同点在于 defer 后跟一个{} 执行该代码块

```
defer {
  // 需要执行的代码，涉及词法，语法分析
}
```
先做个简单的词法分析，明天继续

### 2024-10-10
继续做defer的功能，目前实现是在{}里面的代码支持defer{}操作,比如if 语句中，函数中等
```
fn hello() {
  var_dump("first hello")
  defer {
    var_dump("world")
  }
  var_dump("hello")
}

fn hello_if() {
  let a = 1;
  if (a > 0) {
    defer {
      var_dump("end print")
    }
    var_dump("big than 0")
  }
}

hello()
hello_if()
```

预计输出
```
first hello
hello
world
big than 0
end print
```

### 2024-10-11
继续实现this关键字，这个方法里面使用，指向当前的对象，对于链式的操作很有用，比如
```
this->setName("sevenpan")->setAge(12)->dump();
```
发现一个文件，引入文件的相对路径问题（todo）

### 2024-10-13
实现this在对象方法中的使用，不要词法分析，语法分析，只需要在对应的位置往环境变量注入实例，简单实现


### 2024-10-16
准备实现对象的构造方法，毕竟实际开发中，构造方法还是很实用

```
class User {
  let _name = "seven"
  let _age = 12
  fn __init(name, age) {
    _name = name
    _age = age
  }
  fn dump() {
    var_dump(name)
    var_dump(age)
  }
}
let user = new User("pantingwen", 25)
user->dump()
```
预计输出
```
pantingwen
25
```

思路：
对象的定义里面要支持参数的解析，参考函数
对象初始化之后，检查是否__init 函数，如果有，调用

### 2024-12-13
进一个月没有大的更新，开始继续添加功能，昨天中了一等奖，要加油干活

* 支持了 && 和 || 操作，实际使用的会有多个组合条件判断
* 支持 file_get_contents 和 file_put_contents 来读写文件
* 支持 `__FILE__` 和 `__DIR__` 两个魔法变量，获取当前脚本的信息
* 支持函数的可选默认参数
* 上周支持内置函数 version 获取z语言的版本号
* 支持三元表达式