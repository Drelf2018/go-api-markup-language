from api import get_api, parse_api

API = get_api("path")
# loop

async def demo(args):
    """
    hint
    """
    api = parse_api(API["demo"])
    return await api.update().result
# end