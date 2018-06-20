import os

import requests

PROXY_SERVICE_ADDR = os.environ.get("PROXY_SERVICE_ADDR", "192.168.0.10:12001")


class ProxyProvider:

    def get_proxy(self) -> dict:
        raise NotImplementedError()


class KDProxyProvider(ProxyProvider):

    def __init__(self):
        self.url = "http://{}/proxy/random".format(PROXY_SERVICE_ADDR)
        self.session = requests.session()

    def get_proxy(self) -> dict:

        r = self.session.get(self.url)

        if r.ok:
            try:
                proxy = r.json()

                if proxy.get("login"):
                    return {
                        "http": "http://{}:{}@{}:{}".format(proxy['login'], proxy['password'], proxy['addr'],
                                                            proxy['port']),
                        "https": "http://{}:{}@{}:{}".format(proxy['login'], proxy['password'], proxy['addr'],
                                                              proxy['port'])
                    }
                else:
                    return {
                        "http": "http://{}:{}".format(proxy['addr'], proxy['port']),
                        "https": "http://{}:{}".format(proxy['addr'], proxy['port'])
                    }
            finally:
                r.close()

        else:
            raise Exception("Proxy service not found")



