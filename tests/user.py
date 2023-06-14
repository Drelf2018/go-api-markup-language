from api import get_api, parse_api

API = get_api("./tests/user.json")


async def get_user_card(other: str = "其实也没什么"):
    """
    用户名片信息

    Args:
        other (str): 
    """
    api = parse_api(API["get_user_card"])
    return await api.update(other=other).result
