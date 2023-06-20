from api import get_api, parse_api

API = get_api("./tests/user.json")


async def get_user_card(mid: int, test: int = 0, photo: bool = False):
    """
    用户名片信息

    Args:
        test (num): 

        mid (num): 目标用户mid

        photo (bool): 是否请求用户主页头图
    """
    api = parse_api(API["get_user_card"])
    return await api.update(test=test, mid=mid, photo=photo).result
