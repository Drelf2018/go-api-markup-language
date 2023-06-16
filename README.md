## go-bilibili-api

一种记录 Api 调用方式的文件格式的 golang 解析器。

### 使用

在 `tests` 目录下有文件 `user.aml` 。

```api
GET get_user_card: 用户名片信息 = {
    url = https://api.bilibili.com/x/web-interface/card
    params = {
        int mid: 目标用户mid
        bool photo: 是否请求用户主页头图 = false
    }
}
```

这个文件的格式是：每行文本称为一个 `Token` 每个 `Token` 由四部分组成，分别为 `Type` `Name` `Hint` `Value` 。

以第一行为例 `Type=GET` `Name=get_user_card` `Hint=用户名片信息` `Value={` 。

这些参数只有 `Name` 是必须的，其他可以为空。
例如 `params` 行 `Type=""` `Name=params` `Hint=""` `Value={` 。

运行命令：

```go
go run tests\main.go
```

可以得到 `user.json` `user.py` 文件，其中起主要作用的代码为：

```go
// 读取并解析 user.api
am := parser.GetApi("./tests/user.api")
// 输出为 .json 文件
am.ToJson("./tests/user.json")
// 输出配套的 .py 文件
translator.ToPython(am, "./tests/user.json")
```

于是我们在 `user.py` 中得到这样的代码：

```python
def get_user_card(mid: int, photo: bool = False):
    """
    用户名片信息

    Args:
        photo (bool): 是否请求用户主页头图

        mid (int): 目标用户mid

    """
    api = Api(**API["get_user_card"])
    api.update(photo=photo, mid=mid)
    return api.request()

# 尝试 print(get_user_card(2))
```

这就是自动生成的调用 `user.api` 中接口的 `Python` 代码。

~~之后会支持别的语言大概~~