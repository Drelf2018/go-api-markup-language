from api import get_api, parse_api

API = get_api("./tests/user.json")


async def get_user_card(mid: int, photo: bool = False, test: int = 0):
    """
    用户名片信息

    Args:
        mid (num): 目标用户mid

        photo (bool): 是否请求用户主页头图

        test (num): 
    """
    api = parse_api(API["get_user_card"])
    return await api.update(mid=mid, photo=photo, test=test).result
