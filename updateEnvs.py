import json
import requests
from requests import Response


def conDict(resp: Response) -> dict | list:
    res: dict = json.loads(resp.text)
    return res['data']


def getToken(url: str, client_id: str, client_secret: str) -> str:
    """
    获取青龙面板的token
    :param url: 请求路径
    :param client_id: 密钥id
    :param client_secret: 密钥
    :return: 以字符串形式返回
    """
    resp: Response = requests.get(f'{url}/open/auth/token', params={
        "client_id": client_id,
        "client_secret": client_secret
    })
    data: dict = conDict(resp)
    return f"{data['token_type']} {data['token']}"


def getEnvs(url: str, token: str) -> list:
    """
    获取青龙面板所有的环境变量
    :param url:请求路径
    :param token:青龙面板token
    :return:以字典形式返回
    """
    resp: Response = requests.get(f'{url}/open/envs', headers={
        "Accept": "application/json",
        "Authorization": token
    })
    return conDict(resp)


def updateEnvs(id: int, url: str, token: str, value: str) -> dict:
    """
    更新青龙面板JD_COOKIE环境变量
    :param id:更新环境变量的id
    :param url:请求路径
    :param token:青龙面板token
    :param value:环境变量的value
    :return:
    """
    resp: Response = requests.put(f"{url}/open/envs", data=json.dumps({
        "value": value,
        "name": "JD_COOKIE",
        "remarks": "京东Cookie",
        "id": id
    }), headers={
        "Accept": "application/json",
        "Content-Type": "application/json",
        "Authorization": token
    })
    return json.loads(resp.text)


def main() -> None:
    sourceUrl: str = ''
    sourceClientID: str = ''
    sourceClientSecret: str = ''

    distUrl: str = ''
    distClientID: str = ''
    distClientSecret: str = ''

    global qlId
    sourceToken: str = getToken(
        sourceUrl, sourceClientID, sourceClientSecret)
    cookieDict = getEnvs(sourceUrl, sourceToken)
    cookies: list[str] = [item['value'] for item in cookieDict]
    cookieStr: str = '&'.join(cookies)
    distToken: str = getToken(
        distUrl, distClientID, distClientSecret)
    cookiesInfo: list = getEnvs(distUrl, distToken)

    for item in cookiesInfo:
        if item['name'] == 'JD_COOKIE':
            qlId = item['id']
    res: dict = updateEnvs(qlId, distUrl, distToken, cookieStr)

    if res['code'] == 200:
        print('更新JD_COOKIE环境变量成功!')
    else:
        print('更新JD_COOKIE环境变量失败!')


if __name__ == '__main__':
    main()
