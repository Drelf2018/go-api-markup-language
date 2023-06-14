from api import get_api, parse_api

API = get_api("./tests/user.json")


async def get_user_card():
    """
    用户名片信息
    """
    api = parse_api(API["get_user_card"])
    return await api.update().result
