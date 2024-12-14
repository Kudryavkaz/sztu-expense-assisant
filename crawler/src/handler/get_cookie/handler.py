import json
import base64
import urllib3
import requests

from lxml import etree
from grpc.aio import ServicerContext
from src.grpc.sztuea.v1.crawler_pb2 import GetCookieRequest, GetCookieResponse
from src.log import logger
from src.utils.ocr import ocr

urllib3.disable_warnings()
base_url = 'https://sjdxykt.sztu.edu.cn'
login_url = base_url + '/sso//login'
login_api = base_url + '/sso/doLogin'
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3'}


def get_cookie_handler(request: GetCookieRequest, context: ServicerContext) -> GetCookieResponse:
    logger.info(f"[GetCookie] Account: {request.account}, Password: {request.password}")
    session = requests.session()
    pag_text = session.get(url=login_url, headers=headers).text

    tree = etree.HTML(pag_text, None)

    captcha_code_img = base_url + tree.xpath('//*[@id="captchaCodeImg"]')[0].attrib['src']

    img_content = session.get(captcha_code_img, headers=headers, verify=False).content
    base64_img = base64.b64encode(img_content).decode('utf-8')
    captcha_code = ocr(base64_img)

    data = {
        "loginType": "rftSigner",
        "plat": "",
        "account": request.account,
        "password": request.password,
        "captchaCode": captcha_code,
        "needBind": "",
        "bindPlatform": "",
        "openid": "",
        "unionid": "",
        "alipayUserid": "",
        "ddUserid": "",
        "t": 5,
        "renter": ""
    }

    response = session.post(url=login_api, headers=headers, data=data)

    user_info = json.loads(response.text)
    cookie = user_info.get('msg', None)
    if cookie is None:
        logger.error(f"[GetCookie] Error: {user_info}")
    else:
        logger.info(f"[GetCookie] Cookie: {user_info['msg']}")

    return GetCookieResponse(cookie=cookie)
