import json
from dataclasses import dataclass, field
from typing import Dict

import httpx

SESSION = None


async def get_session() -> httpx.AsyncClient:
    """
    获取异步会话
    """
    global SESSION
    if SESSION is None:
        SESSION = httpx.AsyncClient()
    return SESSION


def new_dict(dic: Dict[str, dict]) -> Dict[str, str]:
    return {k: v.get("value", "") for k, v in dic.items()}


@dataclass
class Api:
    """
    Api 信息类
    """
    url: str
    method: str
    comment: str = ""
    info: Dict[str, str] = field(default_factory=dict)
    data: Dict[str, dict] = field(default_factory=dict)
    params: Dict[str, dict] = field(default_factory=dict)
    headers: Dict[str, dict] = field(default_factory=dict)
    cookies: Dict[str, dict] = field(default_factory=dict)

    def __post_init__(self):
        self.method = self.method.upper()
        self._data = new_dict(self.data.get("value", {}))
        self._params = new_dict(self.params.get("value", {}))
        self._headers = new_dict(self.headers.get("value", {}))
        self._cookies = new_dict(self.cookies.get("value", {}))
        self.__result = None

    @property
    async def result(self):
        """
        获取请求结果
        """
        if self.__result is None:
            session = await get_session()
            resp = await session.request(
                self.method,
                self.url,
                data=self._data,
                params=self._params,
                headers=self._headers,
                cookies=self._cookies,
            )
            self.__result = resp.json()
        return self.__result

    def update_data(self, **kwargs):
        self._data.update(kwargs)
        self.__result = None
        return self

    def update_params(self, **kwargs):
        self._params.update(kwargs)
        self.__result = None
        return self

    def update(self, **kwargs):
        if self.method == "GET":
            return self.update_params(**kwargs)
        else:
            return self.update_data(**kwargs)


KEYS = Api("", "").__dict__.keys()


def parse_api(api: Dict[str, str | dict]):
    """
    解析 api
    """
    info = {}
    for k in list(api.keys()):
        if k not in KEYS:
            info[k] = api.pop(k)
    return Api(info=info, **api)


def get_api(path: str):
    """
    获取 api 列表
    """
    with open(path, "r", encoding="utf-8") as fp:
        return json.load(fp)