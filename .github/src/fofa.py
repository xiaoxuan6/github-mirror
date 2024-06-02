import json
import re
import time

from selenium import webdriver
from webdriver_manager.chrome import ChromeDriverManager


class Fofa:
    def __init__(self):
        self.url = 'https://fofa.info/result?qbase64=IumhueebruWfuuS6jkNsb3VkZmxhcmUgV29ya2Vyc%20' \
                   '%208jOW8gOa6kOS6jkdpdEh1YiIgJiYgY291bnRyeT0iQ04i&page=1&page_size=50'

        chrome_options = webdriver.ChromeOptions()
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--headless')
        chrome_options.add_argument(
            '--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) '
            'Chrome/58.0.3029.110 Safari/537.3')
        # 隐藏正受到自动测试软件控制
        chrome_options.add_experimental_option('excludeSwitches', ['enable-automation'])

        self.driver = webdriver.Chrome(ChromeDriverManager().install(), options=chrome_options)
        self.driver.maximize_window()
        self.driver.get(self.url)
        self.urls = []

    def init(self):
        with open('fofa_cookie.txt') as f:
            cookiesJson = json.loads(f.read())

            for cookie in cookiesJson:
                cookieStr = {
                    "domain": cookie['domain'],
                    "hostOnly": cookie['hostOnly'],
                    "httpOnly": cookie['httpOnly'],
                    "name": cookie['name'],
                    "path": cookie['path'],
                    "secure": cookie['secure'],
                    "session": cookie['session'],
                    "storeId": cookie['storeId'],
                    "value": cookie['value'],
                    "id": cookie['id']
                }

                self.driver.add_cookie(cookie_dict=cookieStr)

        self.driver.refresh()

    def download(self):
        contents = self.driver.find_elements_by_class_name('hsxa-host')
        for content in contents:
            text = re.sub(r'https?://', '', content.text)
            if len(text) > 0:
                self.urls.append(f'{text}')

        if len(self.urls) > 0:
            with open('urls.txt', mode='w', encoding='utf-8') as f:
                f.write(','.join(self.urls))

        self.driver.quit()


if __name__ == '__main__':
    fofa = Fofa()
    fofa.init()
    time.sleep(3)
    fofa.download()
