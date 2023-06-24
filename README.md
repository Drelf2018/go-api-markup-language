# go-bilibili-api

一种记录 Api 调用方式的文件格式的 golang 解析器。

## 介绍

在 `tests` 目录下有文件 `user.aml` 。

```aml
type query = {
    num mid: 目标用户mid
    bool photo: 是否请求用户主页头图 = false
}

type res<T> = {
    num code: 返回码 = 0
    str message: 返回消息
    num ttl = 1
    T data: 数据
}

type card = {
    str mid
    str name
    str face
}

type userInfo<T> = {
    T card: 查询用户的信息
    bool following: 是否关注
    num follower: 粉丝数
}

GET get_user_card: 用户名片信息 = {
    str url = https://api.bilibili.com/x/web-interface/card
    query params
    res<userInfo<card>> response
}
```

这个文件的格式是：每行文本称为一个 `Sentence` 每个 `Sentence` 由四部分组成，分别为 `Type` `Name` `Hint` `Value` 。

```aml
GET get_user_card: 用户名片信息 = {
```

以上面这行为例：`Type=GET` `Name=get_user_card` `Hint=用户名片信息` `Value={` 。

这些参数只有 `Name` 是必须的，其他可以为空。

```aml
res<userInfo<card>> response
```

例如上面这行：`Type=res<userInfo<card>>` `Name=response` `Hint=""` `Value=""` 。

### 关键字

|关键字|意义|
|:----------------|:------------------------------------|
|query|http请求当中的query参数,即url问号后的部分|
|body|http请求当中的请求体参数,一般是json类型,当然也会有二进制,或者表单等格式的数据|
|required|指定参数是必要的|
|optional|指定参数是可选的|
|get|api方法为get请求|
|post|api方法为post请求|
|delete|api方法为delete请求|
|put|api方法为put请求|
|option|api方法为option请求|
|head|api方法为head请求|
|patch|api方法为patch请求|
|enum|字段数据类型为枚举值|
|str|字段数据类型为字符串|
|num|字段数据类型为数字|
|auto|aml会自动推导数据类型,默认为字符串|
|bool|布尔值,即true 或者false|
|import|导入其他文件当中的类型定义|
|from|指定导入来源|
|deprecate|指定参数被弃用|

### 编写接口

```aml
GET get_user_card: 用户名片信息 = {
    str url = https://api.bilibili.com/x/web-interface/card
    query params
    res<userInfo<card>> response
    notice = 使用函数的注意事项
}
```

使用 `GET` 或 `POST` 作为 `Type` 来定义一个 Api 接口，其后的 `get_user_card: 用户名片信息` 是用来导出代码时会用到的函数名和注释。

接口必须使用大括号包裹 `url` 子语句，表示这个接口的地址。其他的子语句是可选的，例如 `params` 和 `response` 。

一些与网络请求相关的字段被称为特殊字段，包含 `url` `data` `params` `headers` `cookies` `response` 。

除此之外都是普通字段，他们也会被导出至 `.json` 或 `.yml` 文件中，但可能不会导出到代码中。

例如上面的 `notice` 字段，注意到该字段并没有显式记录类型，因此内部会根据其初始值判断类型为 `str` 。

```diff
- notice = 使用函数的注意事项
+ auto notice = 使用函数的注意事项
```

你也可以通过设置类型为 `auto` 达到同样的效果。

### 神奇的类型

你可能到了，上述格式中出现了非JSON数据类型（JSON数据类型`str`、`num`、`bool`这一类）

```aml
query params
res<userInfo<card>> response
```

这里的 `query` 和 `res<userInfo<card>>` 是什么鬼啊？

```aml
type query = {
    num mid: 目标用户mid
    bool photo: 是否请求用户主页头图 = false
}
```

找到 `query` 的定义处，我们使用 `type` 关键字定义了这个类型，他的值是一个字典。这样我们就可以在**后续**语句中复用这个类型了。

```diff
- params = {
-     num mid: 目标用户mid
-     bool photo: 是否请求用户主页头图 = false
- }
+ query params
```

你可能也注意到了，上述语句中出现了含有尖括号 `<>` 的类型。没错，这就是泛型。

```aml
type res<T> = {
    num code: 返回码 = 0
    str message: 返回消息
    num ttl = 1
    T data: 数据
}
```

找到`res`的定义处，使用`<>`和其中的任意字符表示这个类型的泛型，泛型可用于其内部语句的类型处。

当`res<userInfo<card>> response`使用时。`response`会被解析为：

```aml
response = {
    num code: 返回码 = 0
    str message: 返回消息
    num ttl = 1
    userInfo<card> data: 数据
}
```

参数可以多个，用 `,` 分隔，也就是 `type res<T1,T2> = xxx` 和 `res<a, b> response`

### 从文件导入类型

```aml
# lib.aml
type query = {
    num mid: 目标用户mid
    bool photo: 是否请求用户主页头图 = false
}

type card = {
    str mid
    str name
    str face
}
```

```aml
# user.aml
from lib import query, card

...

GET get_user_card: 用户名片信息 = {
    str url = https://api.bilibili.com/x/web-interface/card
    query params
    res<userInfo<card>> response
}
```

导入别的文件里的类型就这么简单，然后可以用 `*` 导入全部类型：

```diff
- from lib import query, card
+ from lib import *
```

### 值有什么用

前文提到，每条语句的 `Value` 项并不是必须的，那么写值有什么用呢？目的是方便导出：

```aml
params = {
    num mid: 目标用户mid
    bool photo: 是否请求用户主页头图 = false
}
```

被导出至 `python` 代码时：

```python
async def get_user_card(mid: int, photo: bool = False):
    """
    用户名片信息
    """
    pass  # 以下省略
```

可以发现，值的有无与函数参数中默认值有无是一致的。

如果有一个字段是固定某个值，而我又不希望导出的代码中可以在调用函数时可以修改它，可以在值的后面加上 `,constant`

### 多行文本

当值需要写多行时，可以使用 `"` `'` 等包裹文本：

```aml
GET get_user_card: 用户名片信息 = {
    str url = "https://api.bilibili.com/x/
web-interface/card"
}
```

这会被解析成 `url = https://api.bilibili.com/x/web-interface/card` 。

注意，拼接后的字符串没有换行符。

### 最后

~~之后会支持别的语言大概~~

## 相关项目

- [API Markup Language JavaScript]https://github.com/Kamisato-Ayaka-233/ApiMarkupLanguage) - 使用TypeScript解析AML