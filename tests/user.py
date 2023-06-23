from api import get_api, parse_api

API = get_api("./tests/user.json")


async def get_user_card(mid: int, photo: bool = False):
    """
    用户名片信息

    Args:
        photo (bool): 是否请求用户主页头图

        mid (num): 目标用户mid
    """
    api = parse_api(API["get_user_card"])
    return await api.update(photo=photo, mid=mid).result
