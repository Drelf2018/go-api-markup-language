import json
from dataclasses import dataclass, field
from typing import Dict

import httpx


@dataclass
class Api:
    url: str
    method: str
    comment: str = ""
    data: Dict[str, dict] = field(default_factory=dict)
    params: Dict[str, dict] = field(default_factory=dict)

    def __post_init__(self):
        self.method = self.method.upper()
        self.original_data = self.data.copy()
        self.original_params = self.params.copy()
        self.data = {k: v.get("value", "").replace(",constant", "") for k, v in self.data.items()}
        self.params = {k: v.get("value", "").replace(",constant", "") for k, v in self.params.items()}
        self.__result = None

    def request(self):
        return httpx.request(self.method, self.url, data=self.data, params=self.params).text

    def update_data(self, **kwargs):
        self.data.update(kwargs)
        self.__result = None
        return self

    def update_params(self, **kwargs):
        self.params.update(kwargs)
        self.__result = None
        return self

    def update(self, **kwargs):
        if self.method == "GET":
            return self.update_params(**kwargs)
        else:
            return self.update_data(**kwargs)


def get_api(path: str):
    with open(path, "r", encoding="utf-8") as fp:
        return json.load(fp)


API = get_api("./tests/user.json")


def get_user_card(mid: int, photo: bool = False):
    """
    用户名片信息

    Args:
        mid (int): 目标用户mid

        photo (bool): 是否请求用户主页头图

    """
    api = Api(**API["get_user_card"])
    api.update(mid=mid, photo=photo)
    return api.request()
